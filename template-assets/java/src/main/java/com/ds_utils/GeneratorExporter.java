package com.ds_utils;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.Queue;

import com.fasterxml.jackson.databind.ObjectMapper;

/**
 * Utility class for generating and exporting various data structures.
 */
public class GeneratorExporter {

    /**
     * Generates a binary tree from a list of integer values in BFS order.
     * 
     * @param values the list of integer values
     * @return the root of the generated binary tree
     */
    public static TreeNode generateTree(List<Integer> values) {
        if (values == null || values.isEmpty())
            return null;

        TreeNode root = new TreeNode(values.get(0));
        Queue<TreeNode> queue = new LinkedList<>();
        queue.add(root);

        int i = 1;
        while (!queue.isEmpty() && i < values.size()) {
            TreeNode current = queue.poll();

            if (i < values.size() && values.get(i) != null) {
                current.left = new TreeNode(values.get(i));
                queue.add(current.left);
            } else {
                current.left = null;
            }
            i++;

            if (i < values.size() && values.get(i) != null) {
                current.right = new TreeNode(values.get(i));
                queue.add(current.right);
            } else {
                current.right = null;
            }
            i++;
        }
        return root;
    }

    /**
     * Exports a binary tree to a list of integer values in BFS order.
     * 
     * @param root the root of the binary tree
     * @return the list of integer values representing the tree
     */
    public static List<Integer> exportTree(TreeNode root) {
        List<Integer> result = new ArrayList<>();
        if (root == null)
            return result;

        Queue<TreeNode> queue = new LinkedList<>();
        queue.add(root);

        while (!queue.isEmpty()) {
            TreeNode current = queue.poll();
            if (current != null) {
                result.add(current.val);
                queue.add(current.left);
                queue.add(current.right);
            } else {
                result.add(null);
            }
        }

        // Trim trailing nulls
        while (!result.isEmpty() && result.get(result.size() - 1) == null) {
            result.remove(result.size() - 1);
        }

        return result;
    }

    /**
     * Generates a graph from a list of edges.
     * 
     * @param edges the list of edges, where each edge is represented by a list of two integers
     * @return the generated graph
     * @throws IllegalArgumentException if any edge does not have exactly two elements
     */
    public static Graph generateGraph(List<List<Integer>> edges) {
        Graph graph = new Graph();
        for (List<Integer> edge : edges) {
            if (edge.size() != 2) {
                throw new IllegalArgumentException("Each edge must have exactly two elements");
            }
            graph.addEdge(edge.get(0), edge.get(1));
        }
        return graph;
    }

    /**
     * Exports a graph to a list of edges.
     * 
     * @param graph the graph to be exported
     * @return the list of edges, where each edge is represented by a list of two integers
     */
    public static List<List<Integer>> exportGraph(Graph graph) {
        List<List<Integer>> edges = new ArrayList<>();
        for (Map.Entry<Integer, List<Integer>> entry : graph.adjList.entrySet()) {
            for (Integer neighbor : entry.getValue()) {
                edges.add(Arrays.asList(entry.getKey(), neighbor));
            }
        }
        return edges;
    }

    /**
     * Generates a linked list from a list of integer values.
     * 
     * @param values the list of integer values
     * @return the head of the generated linked list
     */
    public static ListNode generateLinkedList(List<Integer> values) {
        if (values == null || values.isEmpty())
            return null;

        ListNode head = new ListNode(values.get(0));
        ListNode current = head;

        for (int i = 1; i < values.size(); i++) {
            current.next = new ListNode(values.get(i));
            current = current.next;
        }
        return head;
    }

    /**
     * Exports a linked list to a list of integer values.
     * 
     * @param head the head of the linked list
     * @return the list of integer values representing the linked list
     */
    public static List<Integer> exportLinkedList(ListNode head) {
        List<Integer> result = new ArrayList<>();
        ListNode current = head;

        while (current != null) {
            result.add(current.val);
            current = current.next;
        }

        return result;
    }

    private static final ObjectMapper objectMapper = new ObjectMapper();

    /**
     * Exports a list to a JSON string.
     * 
     * @param array the list to be exported
     * @return the JSON string representation of the list
     * @throws IllegalArgumentException if the list cannot be converted to JSON
     */
    public static String exportArray(List<?> array) {
        try {
            return objectMapper.writeValueAsString(array);
        } catch (com.fasterxml.jackson.core.JsonProcessingException e) {
            throw new IllegalArgumentException("Failed to export array", e);
        }

    }

}
