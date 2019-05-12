package com.github.identitii.gdbc.wrapper;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.ResultSet;
import java.sql.ResultSetMetaData;
import java.sql.SQLException;
import java.sql.Statement;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

public class WrappedConnection {
	private ObjectMapper jsonMapper = new ObjectMapper();

	private Connection conn;

	public WrappedConnection(String url, String user, String password) throws SQLException {
		conn = DriverManager.getConnection(url, user, password);
		conn.setAutoCommit(true);
		conn.setTransactionIsolation(Connection.TRANSACTION_SERIALIZABLE);
	}

	/**
	 * Executes an SQL query against the wrapped connection, and returns the ResultSet as JSON.
	 * 
	 * @param query The SQL query to be run
	 * @return The ResultSet serialized as JSON
	 * @throws SQLException
	 */
	public String testQueryJSON(String query) throws SQLException {

		try (Statement stmt = conn.createStatement(); ResultSet rs = stmt.executeQuery(query);) {

			List<Map<String, Object>> rows = new ArrayList<>();
			ResultSetMetaData rsmd = rs.getMetaData();
			int columnCount = rsmd.getColumnCount();

			while (rs.next()) {
				Map<String, Object> row = new HashMap<>();
				for (int i = 1; i <= columnCount; i++) {
					row.put(rsmd.getColumnName(i), rs.getObject(i));
				}
				rows.add(row);
			}
			
			return jsonMapper.writeValueAsString(rows);

		} catch (JsonProcessingException e) {
			try {
				return jsonMapper.writeValueAsString(e.getMessage());
			} catch (JsonProcessingException e1) {return "";} // Can't happen.
		} 

	}

}
