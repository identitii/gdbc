package com.github.identitii.gdbc.wrapper;

import java.util.ArrayList;
import java.util.List;

import org.graalvm.nativeimage.IsolateThread;
import org.graalvm.nativeimage.c.function.CEntryPoint;
import org.graalvm.nativeimage.c.type.CCharPointer;

import static com.github.identitii.gdbc.wrapper.Util.*;

public class DriverWrapper {

	private static List<WrappedConnection> connections = new ArrayList<WrappedConnection>();

	@CEntryPoint(name = "connect")
	public static synchronized int connect(IsolateThread thread, CCharPointer url, CCharPointer user, CCharPointer password) throws Exception {
		connections.add(new WrappedConnection(j(url), j(user), j(password)));
		return connections.size() - 1;
	}

	@CEntryPoint(name = "testQueryJSON")
	public static synchronized CCharPointer testQueryJSON(IsolateThread thread, int connectionID, CCharPointer query) throws Exception {
		return c(connections.get(connectionID).testQueryJSON(j(query)));
	}

}
