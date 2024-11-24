package com.evaluation;
import javax.tools.*;
import java.io.File;
import java.io.StringWriter;
import java.net.URL;
import java.net.URLClassLoader;
import java.util.Collections;

public class JavaCompilerUtil {
    public static Class<?> compileAndLoad(String className, String code) throws Exception {
        // Get the system Java compiler
        JavaCompiler compiler = ToolProvider.getSystemJavaCompiler();
        DiagnosticCollector<JavaFileObject> diagnostics = new DiagnosticCollector<>();

        // Specify the output directory for the compiled class
        File outputDir = new File("target/generated-classes");
        outputDir.mkdirs();

        // Use try-with-resources to ensure the file manager is closed
        try (StandardJavaFileManager fileManager = compiler.getStandardFileManager(diagnostics, null, null)) {
            // Set the output directory
            fileManager.setLocation(StandardLocation.CLASS_OUTPUT, Collections.singletonList(outputDir));

            // Create a JavaFileObject from the provided source code
            JavaFileObject file = new JavaSourceFromString(className, code);
            Iterable<? extends JavaFileObject> compilationUnits = Collections.singletonList(file);

            // Prepare a writer to capture compiler output
            StringWriter writer = new StringWriter();
            // Compile the source code
            boolean success = compiler.getTask(writer, fileManager, diagnostics, null, null, compilationUnits).call();

            // If compilation fails, throw an exception with the compiler output
            if (!success) {
                throw new Exception("Compilation failed: " + writer.toString());
            }
        }

        // Load the compiled class using a URLClassLoader
        try (URLClassLoader classLoader = new URLClassLoader(new URL[]{outputDir.toURI().toURL()})) {
            return classLoader.loadClass(className);
        }
    }
}
