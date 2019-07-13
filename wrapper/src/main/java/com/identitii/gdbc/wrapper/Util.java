package com.identitii.gdbc.wrapper;

import org.graalvm.nativeimage.c.type.CCharPointer;
import org.graalvm.nativeimage.c.type.CTypeConversion;

public class Util {

	public static String j(CCharPointer in) {
		return CTypeConversion.toJavaString(in);
	}

	public static String nullj(CCharPointer in) {
		String out = CTypeConversion.toJavaString(in);
		if (out.equals("")) {
			return null;
		}
		return out;
	}

	public static CCharPointer c(String in) {
		return CTypeConversion.toCString(in).get();
	}

}
