package com.identitii.gdbc.wrapper;

import java.lang.reflect.InvocationHandler;
import java.lang.reflect.Method;
import java.lang.reflect.Proxy;
import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.text.MessageFormat;
import java.util.Arrays;

import org.json.JSONArray;

import oracle.jdbc.dcn.DatabaseChangeRegistration;

class LoggingInvocationHandler implements InvocationHandler {
	
	public static boolean enabled = false;

//	private final static Logger logger = Logger.getLogger(LoggingInvocationHandler.class.getName());
//	static {
//		logger.setLevel(Level.ALL);
//	}
//
//	private static void log(String message, Object... params) {
//		logger.log(Level.FINE, message, params);
//	}
	
	private static MessageFormat format = new MessageFormat(">> {0}.{1}({2}) {3}");

	private static void log(String iface, String method, String args, String result) {
		System.out.println(format.format(new Object[] {iface, method, args, result}));
	}

	private final Object delegate;
	private final Class<?> iface;

	public LoggingInvocationHandler(final Object delegate, final Class<?> iface) {
		this.delegate = delegate;
		this.iface = iface;
	}

	@Override
	public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {

		String argsString = "";
		if (args != null) {
			argsString = new JSONArray(Arrays.asList(args)).toString();
			argsString = argsString.substring(1, argsString.length() - 1);
		}

		try {
			final Object ret = method.invoke(delegate, args);
			if (ret == null) {
				log(iface.getName(), method.getName(), argsString, "");
			} else {
				log(iface.getName(), method.getName(), argsString, ret.toString());
			}
			return ret;
		} catch (Throwable t) {
			t = t.getCause(); // The first exception is always java.lang.reflect.InvocationTargetException which is boring
			log(iface.getName(), method.getName(), argsString, " (EXCEPTION: " + t.getClass().getName() + " - " + t.getMessage() + ")");
			t.printStackTrace();
			throw t;
		}
	}

	public static Connection wrap(Connection obj) {
		if (!enabled) {
			return obj;
		}
		return (Connection) Proxy.newProxyInstance(obj.getClass().getClassLoader(), new Class[] { Connection.class }, new LoggingInvocationHandler(obj, Connection.class));
	}

	public static PreparedStatement wrap(PreparedStatement obj) {
		if (!enabled) {
			return obj;
		}
		return (PreparedStatement) Proxy.newProxyInstance(obj.getClass().getClassLoader(), new Class[] { PreparedStatement.class }, new LoggingInvocationHandler(obj, PreparedStatement.class));
	}

	public static ResultSet wrap(ResultSet obj) {
		if (!enabled) {
			return obj;
		}
		return (ResultSet) Proxy.newProxyInstance(obj.getClass().getClassLoader(), new Class[] { ResultSet.class }, new LoggingInvocationHandler(obj, ResultSet.class));
	}


}