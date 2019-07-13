import java.sql.Connection;
import java.sql.Driver;
import java.sql.ResultSet;
import java.sql.Statement;
import java.util.Properties;

public class SimpleTest {

	public static void main(String[] args) throws Exception {
		System.out.println("starting");
		Driver d = (Driver) Class.forName("org.postgresql.Driver").newInstance();
		System.out.println("got driver");
		Properties info = new Properties();
		info.put("user", "root");
		info.put("password", "password");
		Connection conn = d.connect("jdbc:postgresql://localhost/test?loggerLevel=DEBUG", info);
		System.out.println("connected");

		try (Statement stmt = conn.createStatement(); ResultSet rs = stmt.executeQuery("SELECT 123;");) {
			rs.next();
			System.out.println("result: " + rs.getInt(1));
		}
		System.out.println("done");
	}

}
