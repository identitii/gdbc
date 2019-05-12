package com.github.identitii.gdbc.wrapper;

import org.graalvm.nativeimage.c.type.CCharPointer;
import org.graalvm.nativeimage.c.type.CTypeConversion;

public class Util {

	public static String j(CCharPointer in) {
		return CTypeConversion.toJavaString(in);
	}

	public static CCharPointer c(String in) {
		return CTypeConversion.toCString(in).get();
	}

}
