package com.identitii.gdbc.wrapper;

import static com.identitii.gdbc.wrapper.Util.c;
import static com.identitii.gdbc.wrapper.Util.j;
import static com.identitii.gdbc.wrapper.Util.nullj;

import java.sql.Connection;
import java.sql.Driver;
import java.sql.DriverManager;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.ResultSetMetaData;
import java.sql.SQLException;
import java.sql.Timestamp;
import java.text.MessageFormat;
import java.util.ArrayList;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.graalvm.nativeimage.IsolateThread;
import org.graalvm.nativeimage.c.function.CEntryPoint;
import org.graalvm.nativeimage.c.function.CFunction;
import org.graalvm.nativeimage.c.type.CCharPointer;
import org.json.JSONArray;

public class DriverWrapper {
	
	// // TODO: this log doesn't appear... stderr vs stdout?
	// private final static Logger logger =
	// Logger.getLogger(DriverWrapper.class.getName());
	// static {
	// logger.setLevel(Level.ALL);
	// }
	//
	// private static void log(String message, Object... params) {
	// logger.log(Level.FINE, message, params);
	// }

	protected static void log(String message, Object... params) {
		//if (LoggingInvocationHandler.enabled) {
			System.out.println(MessageFormat.format(message, params));
		//}
	}

	protected static Throwable error = null;
	protected static Connection connection = null;
	protected static List<PreparedStatement> statements = new ArrayList<PreparedStatement>();
	protected static ArrayList<ResultSet> resultSets = new ArrayList<ResultSet>();
	
	/* Import of a CGO function. */
//    @CFunction("on_update")
//    protected static native boolean onUpdate(IsolateThread thread, CCharPointer name, CCharPointer value);

	@CEntryPoint(name = "getError")
	public static CCharPointer getError(IsolateThread thread) {
		if (error == null) {
			return c(null);
		}
		
		try {
			return toError(error);
		} finally {
			error = null;
		}
	}

	public static void setError(Throwable t) {
		error = t;
	}

	public static CCharPointer toError(Throwable t) {
		String message = t.getClass().getName() + ": " + t.getLocalizedMessage();
		if (!message.contains("No suitable driver found")) {
			t.printStackTrace();
		}
		return c(message);
	}

	@CEntryPoint(name = "enableTracing")
	public static void enableTracing(IsolateThread thread, boolean enabled) {
		LoggingInvocationHandler.enabled = enabled;
	}

	@CEntryPoint(name = "openConnection")
	public static CCharPointer openConnection(IsolateThread thread, CCharPointer url, CCharPointer user, CCharPointer password, int txIsolation) {
		
		if (connection != null) {
			return toError(new IllegalStateException("connection has already been opened"));
		}
		
		if (LoggingInvocationHandler.enabled) {
			Enumeration<Driver> drivers = DriverManager.getDrivers();
			while (drivers.hasMoreElements()) {
				log("Driver loaded: {0}", drivers.nextElement().getClass().getName());
			}
		}

		try {
			log("connecting url: {0} user: {1} password: {2}", j(url), nullj(user), nullj(password));
			Connection conn = LoggingInvocationHandler.wrap(DriverManager.getConnection(j(url), nullj(user), nullj(password)));

			conn.setAutoCommit(true);
			conn.setTransactionIsolation(txIsolation);

			connection = conn;
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}
	
	@CEntryPoint(name = "closeConnection")
	public static CCharPointer closeConnection(IsolateThread thread) {
		if (connection != null) {
			try {
				connection.close();
				statements.clear();
				resultSets.clear();
			} catch (Throwable e) {
				return toError(e);
			} finally {
				connection = null;
			}
		}
		return c(null);
	}
	
	@CEntryPoint(name = "isValid")
	public static boolean isValid(IsolateThread thread, int timeout) {
		try {
		
			return connection.isValid(timeout);
		} catch (Throwable e) {
			setError(e);
			return false;
		}
	}
	
	@CEntryPoint(name = "begin")
	public static CCharPointer begin(IsolateThread thread) {
		try {
			connection.setAutoCommit(false);
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}
	
	@CEntryPoint(name = "commit")
	public static CCharPointer commit(IsolateThread thread) {
		try {
			connection.commit();
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}
	
	@CEntryPoint(name = "rollback")
	public static CCharPointer rollback(IsolateThread thread) {
		try {
			connection.rollback();
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}

	@CEntryPoint(name = "prepare")
	public static int prepare(IsolateThread thread, CCharPointer sql) {
		try {
			PreparedStatement s = LoggingInvocationHandler.wrap(connection.prepareStatement(j(sql)));
			statements.add(s);
			return statements.size() - 1;
		} catch (Throwable e) {
			setError(e);
			return -1;
		}
	}

	@CEntryPoint(name = "closeStatement")
	public static CCharPointer closeStatement(IsolateThread thread, int statement) {
		try {
			if (statements.get(statement) != null) {
				statements.get(statement).close();
				statements.set(statement, null);
			}
			ensureSize(resultSets, statement+1);
			resultSets.set(statement, null);
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}

	@CEntryPoint(name = "numInput")
	public static int numInput(IsolateThread thread, int statement) {
		return -1;/*
		try {
			return statements.get(statement).getParameterMetaData().getParameterCount();
		} catch (Throwable e) {
			setError(e);
			return -1; // Means ignore me in the go api
		}*/
	}
	
	@CEntryPoint(name = "execute")
	public static int execute(IsolateThread thread, int statement) {
		try {
			boolean hasResult = statements.get(statement).execute();
			if (hasResult) {
				throw new SQLException("unexpected results on statement execute");
			}
			return statements.get(statement).getUpdateCount();
		} catch (Throwable e) {
			setError(e);
			return -1;
		} finally {
			try {
				statements.get(statement).close();
			} catch(Throwable t) {
				// Ignored. TODO: Do I have to handle this?
			}
		}
	}
	
	@CEntryPoint(name = "query")
	public static boolean query(IsolateThread thread, int statement) {
		
		try {
			PreparedStatement stmt = statements.get(statement);
			boolean hasResult = stmt.execute();
			ensureSize(resultSets, statement+1);
			resultSets.set(statement, LoggingInvocationHandler.wrap(stmt.getResultSet())); 
			return hasResult;
		} catch (Throwable e) {
			setError(e);
			return false;
		}
	}

	private static void ensureSize(ArrayList<?> list, int size) {
		list.ensureCapacity(size);
		while (list.size() < size) {
			list.add(null);
		}
	}
	
	@CEntryPoint(name = "getMoreResults")
	public static boolean getMoreResults(IsolateThread thread, int statement) {
		try {
			return statements.get(statement).getMoreResults();
		} catch (Throwable e) {
			setError(e);
			return false;
		}
	}
	
	@CEntryPoint(name = "nextResultSet")
	public static boolean nextResultSet(IsolateThread thread, int statement) {
		try {
			ensureSize(resultSets, statement+1);
			resultSets.set(statement, LoggingInvocationHandler.wrap(statements.get(statement).getResultSet()));
			return true;
		} catch (Throwable e) {
			setError(e);
			return false;
		}
	}
	
	@CEntryPoint(name = "columns")
	public static CCharPointer columns(IsolateThread thread, int statement) {
		try {
			ResultSetMetaData md = resultSets.get(statement).getMetaData();
			
			List<String> names = new ArrayList<String>();
			List<String> types = new ArrayList<String>();
			
			for (int i = 1; i <= md.getColumnCount(); i++) { // Columns start at 1
				names.add(md.getColumnName(i));
				types.add(md.getColumnClassName(i));
			}
			
			return c(String.join(",", names) + "|" + String.join(",", types));
			
		} catch (Throwable e) {
			setError(e);
			return c(null);
		}
	}
	
	@CEntryPoint(name = "next")
	public static boolean next(IsolateThread thread, int statement) {
		try {
			return resultSets.get(statement).next();
		} catch (Throwable e) {
			setError(e);
			return false;
		}
	}
	
	@CEntryPoint(name = "setByte")
	public static CCharPointer setByte(IsolateThread thread, int statement, int index, byte value) {
		try {
			statements.get(statement).setByte(index, value);
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}
	
	@CEntryPoint(name = "getByte")
	public static byte getByte(IsolateThread thread, int statement, int index) {
	    try {
	        return resultSets.get(statement).getByte(index);
	    } catch (Throwable e) {
	        setError(e);
	        return -1;
	    }
	}

	@CEntryPoint(name = "setShort")
	public static CCharPointer setShort(IsolateThread thread, int statement, int index, short value) {
		try {
			statements.get(statement).setShort(index, value);
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}
	
	@CEntryPoint(name = "getShort")
	public static short getShort(IsolateThread thread, int statement, int index) {
		try {
			return resultSets.get(statement).getShort(index);
		} catch (Throwable e) {
            setError(e);
            return -1;
		}
	}

	@CEntryPoint(name = "setInt")
	public static CCharPointer setInt(IsolateThread thread, int statement, int index, int value) {
		try {
			statements.get(statement).setInt(index, value);
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}
	
	@CEntryPoint(name = "getInt")
	public static int getInt(IsolateThread thread, int statement, int index) {
		try {
			return resultSets.get(statement).getInt(index);
		} catch (Throwable e) {
            setError(e);
            return -1;
		}
	}

	@CEntryPoint(name = "setLong")
	public static CCharPointer setLong(IsolateThread thread, int statement, int index, long value) {
		try {
			statements.get(statement).setLong(index, value);
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}
	
	@CEntryPoint(name = "getLong")
	public static long getLong(IsolateThread thread, int statement, int index) {
		try {
			return resultSets.get(statement).getLong(index);
		} catch (Throwable e) {
            setError(e);
            return -1;
		}
	}

	@CEntryPoint(name = "setFloat")
	public static CCharPointer setFloat(IsolateThread thread, int statement, int index, float value) {
		try {
			statements.get(statement).setFloat(index, value);
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}
	
	@CEntryPoint(name = "getFloat")
	public static float getFloat(IsolateThread thread, int statement, int index) {
		try {
			return resultSets.get(statement).getFloat(index);
		} catch (Throwable e) {
            setError(e);
            return -1;
		}
	}

	@CEntryPoint(name = "setDouble")
	public static CCharPointer setDouble(IsolateThread thread, int statement, int index, double value) {
		try {
			statements.get(statement).setDouble(index, value);
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}
	
	@CEntryPoint(name = "getDouble")
	public static double getDouble(IsolateThread thread, int statement, int index) {
		try {
			return resultSets.get(statement).getDouble(index);
		} catch (Throwable e) {
            setError(e);
            return -1;
		}
	}
	
	@CEntryPoint(name = "getBigDecimal")
	public static CCharPointer getBigDecimal(IsolateThread thread, int statement, int index) {
		try {
			String value = resultSets.get(statement).getBigDecimal(index).toString();
			return c(value);
		} catch (Throwable e) {
            setError(e);
            return c(null);
		}
	}

	@CEntryPoint(name = "setString")
	public static CCharPointer setString(IsolateThread thread, int statement, int index, CCharPointer value) {
		try {
			statements.get(statement).setString(index, j(value));
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}
	
	@CEntryPoint(name = "getString")
	public static CCharPointer getString(IsolateThread thread, int statement, int index) {
		try {
			return c(resultSets.get(statement).getString(index));
		} catch (Throwable e) {
			setError(e);
			return c(null);
		}
	}

	@CEntryPoint(name = "setTimestamp")
	public static CCharPointer setTimestamp(IsolateThread thread, int statement, int index, long value) {
		try {
			statements.get(statement).setTimestamp(index, new Timestamp(value));
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}
	
	@CEntryPoint(name = "getTimestamp")
	public static long getTimestamp(IsolateThread thread, int statement, int index) {
		try {
			return resultSets.get(statement).getTimestamp(index).getTime();
		} catch (Throwable e) {
			setError(e);
			return -1;
		}
	}

	@CEntryPoint(name = "setNull")
	public static CCharPointer setNull(IsolateThread thread, int statement, int index) {
		try {
			statements.get(statement).setObject(index, null);
			return c(null);
		} catch (Throwable e) {
			return toError(e);
		}
	}
 
	@CEntryPoint(name = "testQueryJSON")
	public static CCharPointer testQueryJSON(IsolateThread thread, CCharPointer query) {
		log("testQueryJSON(query:{1})", j(query));

		int statementId = prepare(thread, query);

		try (PreparedStatement stmt = statements.get(statementId); ResultSet rs = LoggingInvocationHandler.wrap(stmt.executeQuery());) {
			
			JSONArray rows = new JSONArray();
			ResultSetMetaData rsmd = rs.getMetaData();
			int columnCount = rsmd.getColumnCount();

			while (rs.next()) {
				Map<String, Object> row = new HashMap<>();
				for (int i = 1; i <= columnCount; i++) {
					row.put(rsmd.getColumnName(i), rs.getObject(i));
				}
				rows.put(row);
			}

			return c(rows.toString());

		} catch (Throwable e) {
			setError(e);
			return c(null);
		}
	}

}
