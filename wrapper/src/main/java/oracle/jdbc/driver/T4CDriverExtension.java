package oracle.jdbc.driver;

import java.sql.Connection;
import java.sql.SQLException;
import java.util.Properties;

// GDBC NOTE: The only thing changed in this file is to make the class public.

public class T4CDriverExtension extends OracleDriverExtension {
	private static final String _Copyright_2007_Oracle_All_Rights_Reserved_ = null;
	public static final String BUILD_DATE = "Thu_Apr_04_15:06:58_PDT_2013";
	public static final boolean TRACE = false;

	Connection getConnection(String var1, Properties var2) throws SQLException {
		return new T4CConnection(var1, var2, this);
	}

	OracleStatement allocateStatement(PhysicalConnection var1, int var2, int var3) throws SQLException {
		return new T4CStatement(var1, var2, var3);
	}

	OraclePreparedStatement allocatePreparedStatement(PhysicalConnection var1, String var2, int var3, int var4)
			throws SQLException {
		return new T4CPreparedStatement(var1, var2, var3, var4);
	}

	OracleCallableStatement allocateCallableStatement(PhysicalConnection var1, String var2, int var3, int var4)
			throws SQLException {
		return new T4CCallableStatement(var1, var2, var3, var4);
	}

	OracleInputStream createInputStream(OracleStatement var1, int var2, Accessor var3) throws SQLException {
		return new T4CInputStream(var1, var2, var3);
	}
}