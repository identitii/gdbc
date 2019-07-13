package com.identitii.gdbc.wrapper;

import static com.identitii.gdbc.wrapper.Util.c;
import static com.identitii.gdbc.wrapper.Util.j;

import java.lang.reflect.Proxy;
import java.util.ArrayList;
import java.util.List;
import java.util.Properties;

import org.graalvm.nativeimage.IsolateThread;
import org.graalvm.nativeimage.c.function.CEntryPoint;
import org.graalvm.nativeimage.c.function.CFunction;
import org.graalvm.nativeimage.c.type.CCharPointer;
import org.json.JSONObject;

import com.alibaba.fastjson.JSON;

import oracle.jdbc.OracleConnection;
import oracle.jdbc.OracleStatement;
import oracle.jdbc.dcn.DatabaseChangeEvent;
import oracle.jdbc.dcn.DatabaseChangeListener;
import oracle.jdbc.dcn.DatabaseChangeRegistration;
import oracle.jdbc.driver.OracleDriver;

public class OracleExtensions {

	private static List<DatabaseChangeRegistration> registrations = new ArrayList<DatabaseChangeRegistration>();

	/* Import of a CGO function. */  
	@CFunction("oracle_on_change_event")
	protected static native boolean OnChangeEvent(long regID, CCharPointer event);

	@CEntryPoint(name = "oracleRegisterDatabaseChangeNotification")
	public static int registerDatabaseChangeNotification(IsolateThread thread, CCharPointer coptions) {
		OnChangeEvent(123, c("starting registration"));
		DriverWrapper.log("oracleRegisterDatabaseChangeNotification {0}", j(coptions));
		try {
			OracleConnection conn = (OracleConnection) DriverWrapper.connection;

			Properties options = new Properties();

			JSONObject jsonOptions = new JSONObject(j(coptions));
			for (String key : jsonOptions.keySet()) {
				options.put(key, jsonOptions.get(key)+"");
			}

			DatabaseChangeRegistration r = conn.registerDatabaseChangeNotification(options);
			
			r.addListener(new DatabaseChangeListener() {
				@Override
				public void onDatabaseChangeNotification(DatabaseChangeEvent event) {
					DriverWrapper.log("EVENT! {0}", event.getRegId());
					String eventString = JSON.toJSONString(event);
					DriverWrapper.log("EVENT! {0} {1}", event.getRegId(), eventString);
					
					OnChangeEvent(event.getRegId(), c(eventString));
				}
			});

			registrations.add(r);
			return registrations.size() - 1;
		} catch (Throwable e) {
			DriverWrapper.setError(e);
			return -1;
		}
	}

	@CEntryPoint(name = "oracleRegistrationGetRegId")
	public static long registrationGetRegId(IsolateThread thread, int registration) {
		try {
			return registrations.get(registration).getRegId();
		} catch (Throwable e) {
			DriverWrapper.setError(e);
			return -1;
		}
	}
	
	@CEntryPoint(name = "oracleRegistrationGetTables")
	public static CCharPointer registrationGetTables(IsolateThread thread, int registration) {
		try {
			return c(String.join(",", registrations.get(registration).getTables()));
		} catch (Throwable e) {
			DriverWrapper.setError(e);
			return c(null);
		}
	}
	
	@CEntryPoint(name = "oracleRegistrationGetState")
	public static CCharPointer registrationGetState(IsolateThread thread, int registration) {
		try {
			return c(registrations.get(registration).getState().toString());
		} catch (Throwable e) {
			DriverWrapper.setError(e);
			return c(null);
		}
	}

	@CEntryPoint(name = "oracleSetDatabaseChangeRegistration")
	public static CCharPointer setDatabaseChangeRegistration(IsolateThread thread, int statement, int registration) {
		try {
			DatabaseChangeRegistration dcr = registrations.get(registration);
			((OracleStatement)DriverWrapper.statements.get(statement)).setDatabaseChangeRegistration(dcr);
			return c(null);
		} catch (Throwable e) {
			return DriverWrapper.toError(e);
		}
	}

	@CEntryPoint(name = "oracleDriverBuildDate")
	public static CCharPointer driverBuildDate(IsolateThread thread) {
		return c(OracleDriver.BUILD_DATE);
	}

}
