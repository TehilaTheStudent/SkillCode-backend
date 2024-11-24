package com.evaluation;

import java.net.URI;

import javax.tools.SimpleJavaFileObject;

/**
 * A file object used to represent source coming from a string.
 * This class extends SimpleJavaFileObject and provides the source code
 * as a CharSequence.
 */
public class JavaSourceFromString extends SimpleJavaFileObject {
    private final String code;

    /**
     * Constructs a new JavaSourceFromString.
     *
     * @param name the name of the class
     * @param code the source code of the class
     */
    public JavaSourceFromString(String name, String code) {
        // Create a URI for the source code with the given class name
        super(URI.create("string:///" + name.replace('.', '/') + Kind.SOURCE.extension), Kind.SOURCE);
        this.code = code;
    }

    /**
     * Returns the source code of the class.
     *
     * @param ignoreEncodingErrors if true, encoding errors are ignored
     * @return the source code as a CharSequence
     */
    @Override
    public CharSequence getCharContent(boolean ignoreEncodingErrors) {
        return code;
    }
}
