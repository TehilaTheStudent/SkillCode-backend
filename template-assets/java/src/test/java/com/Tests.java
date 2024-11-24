package com;

import java.util.Arrays;
import java.util.List;

import static org.testng.Assert.assertEquals;
import org.testng.annotations.Test;

import com.ds_utils.AbstractType;
import com.ds_utils.Graph;
import com.ds_utils.ListNode;
import com.ds_utils.TreeNode;
import com.ds_utils.TypeConverter;
import com.evaluation.FunctionConfig;
import com.evaluation.TestCase;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.PropertyNamingStrategies;

public class Tests {

    @Test
    public void testMain() {
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

    @Test
    public void testListyToType_IntegerArray() {
        // Input JSON string representing a list of integers
        String jsonInput = "[1, 1, 1]";

        // Define the abstract type for an Array of Integers
        AbstractType integerType = new AbstractType("Integer", null);

        AbstractType arrayType = new AbstractType("Array", integerType);

        // Convert JSON input to the desired type
        Object converted = TypeConverter.listyToType(jsonInput, arrayType);

        // Verify the result is a List of Integers
        List<Integer> expected = Arrays.asList(1, 1, 1);
        assertEquals(converted, expected, "The converted object does not match the expected List<Integer>");
    }

    @Test
    public void testListyToType_Matrix() {
        String jsonInput = "[[1, 2], [3, 4]]";
        AbstractType integerType = new AbstractType("Integer", null);
        AbstractType matrixType = new AbstractType("Matrix", integerType);
        Object converted = TypeConverter.listyToType(jsonInput, matrixType);
        List<List<Integer>> expected = Arrays.asList(Arrays.asList(1, 2), Arrays.asList(3, 4));
        assertEquals(converted, expected, "The converted object does not match the expected Matrix<List<Integer>>");
    }

    @Test
    public void testListyToType_TreeNode() {
        String jsonInput = "[1, 3, 2]";
        AbstractType integerType = new AbstractType("Integer", null);
        AbstractType treeNodeType = new AbstractType("TreeNode", integerType);
        Object converted = TypeConverter.listyToType(jsonInput, treeNodeType);
        TreeNode expected = new TreeNode(1);
        expected.right = new TreeNode(3);
        expected.left = new TreeNode(2);
        assertEquals(converted, expected, "The converted object does not match the expected TreeNode<Integer>");
    }

    @Test
    public void testListyToType_ListNode() {
        String jsonInput = "[1, 2, 3]";
        AbstractType integerType = new AbstractType("Integer", null);
        AbstractType listNodeType = new AbstractType("ListNode", integerType);
        Object converted = TypeConverter.listyToType(jsonInput, listNodeType);
        ListNode expected = new ListNode(1);
        expected.next = new ListNode(2);
        expected.next.next = new ListNode(3);
        assertEquals(converted, expected, "The converted object does not match the expected ListNode<Integer>");
    }

    @Test
    public void testListyToType_Graph() {
        String jsonInput = "[[0, 1], [1, 2], [2, 0]]";
        AbstractType integerType = new AbstractType("Integer", null);
        AbstractType graphType = new AbstractType("Graph", integerType);
        Object converted = TypeConverter.listyToType(jsonInput, graphType);
        Graph expected = new Graph();
        expected.adjList.put(0, Arrays.asList(1));
        expected.adjList.put(1, Arrays.asList(2));
        expected.adjList.put(2, Arrays.asList(0));
        assertEquals(converted, expected, "The converted object does not match the expected Graph<Integer>");
    }
}
