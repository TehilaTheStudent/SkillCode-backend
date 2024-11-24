package com.ds_utils;

import java.util.Arrays;
import java.util.HashMap; // For JSON parsing
import java.util.List; // For List, Map, HashMap, Arrays, and other utility classes
import java.util.Map; // For stream and collect operations
import java.util.stream.Collectors;

import com.fasterxml.jackson.databind.ObjectMapper;

/**
 * Utility class for converting JSON string representations of various data structures
 * into their corresponding Java object representations.
 */
public class TypeConverter {

    private static final ObjectMapper objectMapper = new ObjectMapper();

    /**
     * Converts a JSON string representation of a listy structure into a Java object
     * based on the provided abstract type.
     *
     * @param stringyListyRep the JSON string representation of the listy structure
     * @param abstractType    the abstract type describing the structure
     * @return the Java object representation of the listy structure
     * @throws IllegalArgumentException if the JSON parsing fails
     */
    public static Object listyToType(String stringyListyRep, AbstractType abstractType) {
        try {
            // Parse the JSON string representation into a Java Object
            Object listyRep = objectMapper.readValue(stringyListyRep, Object.class);
            return convertListyToType(listyRep, abstractType);
        } catch (com.fasterxml.jackson.core.JsonProcessingException e) {
            throw new IllegalArgumentException("Failed to parse JSON input: " + stringyListyRep, e);
        }
    }

    /**
     * Recursively converts a parsed listy representation into a Java object
     * based on the provided abstract type.
     *
     * @param listyRep     the parsed listy representation
     * @param abstractType the abstract type describing the structure
     * @return the Java object representation of the listy structure
     */
    private static Object convertListyToType(Object listyRep, AbstractType abstractType) {
        String baseType = abstractType.type;
        AbstractType typeChildren = abstractType.typeChildren;

        // Handle atomic types
        if (Arrays.asList("Integer", "Boolean", "String", "Double").contains(baseType)) {
            return listyRep; // Directly return atomic values
        }

        // Handle composite types
        if ("Array".equals(baseType)) {
            // Convert listyRep to a List of Integers and recursively convert each item
            List<Integer> list = objectMapper.convertValue(listyRep,
                    objectMapper.getTypeFactory().constructCollectionType(List.class, Integer.class));
            return list.stream()
                    .map(item -> convertListyToType(item, typeChildren))
                    .collect(Collectors.toList());
        }

        if ("Matrix".equals(baseType)) {
            // Convert listyRep to a List of Lists of Integers and recursively convert each row
            List<List<Integer>> matrix = objectMapper.convertValue(listyRep,
                    objectMapper.getTypeFactory().constructCollectionType(List.class,
                            objectMapper.getTypeFactory().constructCollectionType(List.class, Integer.class)));
            return matrix.stream()
                    .map(row -> objectMapper.convertValue(row,
                            objectMapper.getTypeFactory().constructCollectionType(List.class, Integer.class)))
                    .collect(Collectors.toList());
        }

        if ("TreeNode".equals(baseType)) {
            // Convert listyRep to a TreeNode structure
            if (!"Integer".equals(typeChildren.type)) {
                throw new IllegalArgumentException("TreeNode can only be of type Integer");
            }
            List<Integer> listyRepAsList = objectMapper.convertValue(listyRep,
                    objectMapper.getTypeFactory().constructCollectionType(List.class, Integer.class));
            TreeNode tree = GeneratorExporter.generateTree(listyRepAsList);
            tree.val = (Integer) convertListyToType(tree.val, typeChildren);

            // Recursively convert left and right subtrees
            if (tree.left != null) {
                Object leftSubtree = convertListyToType(
                        GeneratorExporter.exportTree(tree.left),
                        typeChildren);
                if (leftSubtree instanceof TreeNode) {
                    tree.left = (TreeNode) leftSubtree;
                } else {
                    throw new IllegalArgumentException("Invalid type for left subtree: "
                            + (leftSubtree != null ? leftSubtree.getClass() : "null"));
                }
            }

            if (tree.right != null) {
                Object rightSubtree = convertListyToType(
                        GeneratorExporter.exportTree(tree.right),
                        typeChildren);
                if (rightSubtree instanceof TreeNode) {
                    tree.right = (TreeNode) rightSubtree;
                } else {
                    throw new IllegalArgumentException(
                            "Invalid type for right subtree: "
                                    + (rightSubtree != null ? rightSubtree.getClass() : "null"));
                }
            }

            return tree;
        }

        if ("ListNode".equals(baseType)) {
            // Convert listyRep to a ListNode structure
            if (!"Integer".equals(typeChildren.type)) {
                throw new IllegalArgumentException("ListNode can only be of type Integer");
            }
            List<Integer> listyRepAsList = objectMapper.convertValue(listyRep,
                    objectMapper.getTypeFactory().constructCollectionType(List.class, Integer.class));
            ListNode linkedList = GeneratorExporter.generateLinkedList(listyRepAsList);
            ListNode current = linkedList;
            while (current != null) {
                current.val = (Integer) convertListyToType(current.val, typeChildren);
                current = current.next;
            }
            return linkedList;
        }

        if ("Graph".equals(baseType)) {
            // Convert listyRep to a Graph structure
            if (!"Integer".equals(typeChildren.type)) {
                throw new IllegalArgumentException("Graph can only be of type Integer");
            }
            List<List<Integer>> listyRepAsList = objectMapper.convertValue(listyRep,
                    objectMapper.getTypeFactory().constructCollectionType(List.class,
                            objectMapper.getTypeFactory().constructCollectionType(List.class, Integer.class)));
            Graph graph = GeneratorExporter.generateGraph(listyRepAsList);
            Map<Integer, List<Integer>> newAdjList = new HashMap<>();
            for (Map.Entry<Integer, List<Integer>> entry : graph.adjList.entrySet()) {
                Integer convertedU = (Integer) convertListyToType(entry.getKey(), typeChildren);
                List<Integer> convertedNeighbors = entry.getValue().stream()
                        .map(neighbor -> (Integer) convertListyToType(neighbor, typeChildren))
                        .collect(Collectors.toList());
                newAdjList.put(convertedU, convertedNeighbors);
            }
            graph.adjList = newAdjList;
            return graph;
        }

        throw new IllegalArgumentException("Unsupported type: " + baseType);
    }
}
