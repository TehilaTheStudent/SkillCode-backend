package com.evaluation;

import java.lang.reflect.Method;
import java.util.ArrayList;
import java.util.List;

import com.ds_utils.AbstractType;
import com.ds_utils.TypeConverter;

/**
 * Class for evaluating user-provided code against a set of test cases.
 */
public class CodeEvaluator {

    /**
     * Evaluates the user-provided code against the given test cases.
     *
     * @param userCode       the Java source code provided by the user
     * @param testCases      the list of test cases to evaluate the code against
     * @param functionConfig the configuration of the function to be tested
     * @return a string indicating whether the tests passed or failed
     */
    public static String evaluateUserCode(String userCode, List<TestCase> testCases, FunctionConfig functionConfig) {
        try {
            // Compile and load the user's code
            Class<?> userClass = JavaCompilerUtil.compileAndLoad("UserSolution", userCode);
            for (Method method : userClass.getDeclaredMethods()) {
                // System.out.println("Available Method: " + method.getName() + " with parameters: " +
                //         Arrays.toString(method.getParameterTypes()));
            }

            Object userInstance = userClass.getDeclaredConstructor().newInstance();

            // Get the method to be tested from the compiled class
            Method userMethod;
            try {
                 userMethod = userClass.getMethod(functionConfig.functionName,
                        functionConfig.parameters.stream()
                                .map(AbstractType::getJavaType)
                                .toArray(Class<?>[]::new));
                // System.out.println("Found method: " + userMethod);
            } catch (NoSuchMethodException e) {
                // System.out.println("Method not found: " + e.getMessage());
                throw e;
            }

            // Iterate over each test case and evaluate the user's code
            for (TestCase testCase : testCases) {
                // Convert the test case inputs to the appropriate types
                List<Object> inputs = convertInputs(testCase.parameters, functionConfig.parameters);
                // Convert the expected output to the appropriate type
                Object expected = TypeConverter.listyToType(testCase.expectedOutput, functionConfig.returnType);
                // Invoke the user's method with the test case inputs
                Object result = userMethod.invoke(userInstance, inputs.toArray());

                // Check if the result matches the expected output
                if (!result.equals(expected)) {
                    return "Test failed!";
                }
            }

            return "All tests passed!";
        } catch (Exception e) {
            return "Error during evaluation: " + e.getMessage();
        }
    }

    /**
     * Converts the test case input parameters to the appropriate types.
     *
     * @param parameters the list of input parameters as strings
     * @param types      the list of expected types for the input parameters
     * @return a list of input parameters converted to the appropriate types
     */
    private static List<Object> convertInputs(List<String> parameters, List<AbstractType> types) {
        List<Object> inputs = new ArrayList<>();
        for (int i = 0; i < parameters.size(); i++) {
            inputs.add(TypeConverter.listyToType(parameters.get(i), types.get(i)));
        }
        return inputs;
    }
}
