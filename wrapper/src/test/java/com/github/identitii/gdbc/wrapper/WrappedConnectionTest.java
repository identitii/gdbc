package com.github.identitii.gdbc.wrapper;

import static org.junit.Assert.*;

import java.sql.SQLException;

import org.junit.BeforeClass;
import org.junit.Test;

public class WrappedConnectionTest {

	private static WrappedConnection connection;
	
	@BeforeClass
	public static void connect() throws Exception {
		connection = new WrappedConnection("jdbc:mysql://localhost/test","root","password");
	}
	
	@Test
	public void testTestQueryJSON() throws SQLException {
		String result = connection.testQueryJSON("SELECT 1");
		System.out.println("testTestQueryJSON: " + result);
	}

}
