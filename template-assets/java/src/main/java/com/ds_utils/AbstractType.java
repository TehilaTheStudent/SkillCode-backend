package com.ds_utils;

import java.util.List;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true) // Ignore properties not defined in the class
public class AbstractType {
   public String type;
    public AbstractType typeChildren;

    // Default constructor for Jackson
    public AbstractType() {}

    // Constructor with parameters
    @JsonCreator
    public AbstractType(@JsonProperty("type") String type,
                        @JsonProperty("type_children") AbstractType typeChildren) {
        this.type = type;
        this.typeChildren = typeChildren;
    }

    @Override
    public String toString() {
        return "AbstractType{" +
                "type='" + type + '\'' +
                ", typeChildren=" + typeChildren +
                '}';
    }

    /**
     * Converts this AbstractType to the corresponding Java Class<?>.
     *
     * @return the corresponding Java Class<?> for this type.
     */
    public Class<?> getJavaType() {
        switch (type) {
            case "Integer":
                return Integer.class;
            case "Double":
                return Double.class;
            case "String":
                return String.class;
            case "Boolean":
                return Boolean.class;
            case "Array":
                return List.class; // Arrays can be represented as List in Java
            case "Matrix":
                return List.class; // Matrix will be represented as List<List>
            case "TreeNode":
                return TreeNode.class;
            case "ListNode":
                return ListNode.class;
            case "Graph":
                return Graph.class;
            default:
                throw new IllegalArgumentException("Unsupported type: " + type);
        }
    }
}
