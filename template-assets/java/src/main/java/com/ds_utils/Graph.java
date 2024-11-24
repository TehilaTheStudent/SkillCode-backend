package com.ds_utils;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
public class Graph {
    public Map<Integer, List<Integer>> adjList = new HashMap<>();

    public void addEdge(Integer u, Integer v) {
        adjList.computeIfAbsent(u, k -> new ArrayList<>()).add(v);
    }

    @Override
    public boolean equals(Object obj) {
        if (!(obj instanceof Graph)) return false;
        Graph other = (Graph) obj;
        return Objects.equals(adjList, other.adjList);
    }

    @Override
    public int hashCode() {
        return Objects.hash(adjList);
    }

    @Override
    public String toString() {
        return "Graph{" + "adjList=" + adjList + '}';
    }
}
