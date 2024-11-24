package com.generated;

import java.util.List;

import com.evaluation.FunctionConfig;
import com.evaluation.TestCase;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.PropertyNamingStrategies;

public class Main {
        public static void main(String[] args) {
                try {
                        System.out.println(System.getProperty("java.class.path"));

                        // JSON strings for test cases and function configuration
                        String testCasesJson = "[{\"parameters\": [\"[7,6,4,3,1]\"], \"expected_output\": \"0\"}]";
                        String functionConfigJson = "{\n" +
                                        "  \"name\": \"maxProfit\",\n" +
                                        "  \"parameters\": [\n" +
                                        "    {\n" +
                                        "      \"name\": \"prices\",\n" +
                                        "      \"param_type\": {\n" +
                                        "        \"type\": \"Array\",\n" +
                                        "        \"type_children\": {\n" +
                                        "          \"type\": \"Integer\"\n" +
                                        "        }\n" +
                                        "      }\n" +
                                        "    }\n" +
                                        "  ],\n" +
                                        "  \"return_type\": {\n" +
                                        "    \"type\": \"Integer\"\n" +
                                        "  }\n" +
                                        "}";
                        ObjectMapper mapper = new ObjectMapper();
                        mapper.setPropertyNamingStrategy(PropertyNamingStrategies.SNAKE_CASE);

                        // Deserialize JSON strings into Java objects
                        List<TestCase> testCases = mapper.readValue(
                                        testCasesJson,
                                        mapper.getTypeFactory().constructCollectionType(List.class, TestCase.class));
                        FunctionConfig functionConfig = mapper.readValue(functionConfigJson, FunctionConfig.class);

                        // Print parsed objects to verify correctness
                        System.out.println("Parsed Test Cases: " + testCases);
                        System.out.println("Parsed Function Config: " + functionConfig);

                } catch (Exception e) {
                        e.printStackTrace();
                }
                // Evaluate the user code (assuming evaluateUserCode works correctly)
                // String message = CodeEvaluator.evaluateUserCode(userCode, testCases,
                // functionConfig);
                // System.out.println(message);
        }
}
