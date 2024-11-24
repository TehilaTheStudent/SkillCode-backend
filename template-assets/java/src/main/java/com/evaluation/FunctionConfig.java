package com.evaluation;

import java.util.List;

import com.ds_utils.AbstractType;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

public class FunctionConfig {
    @JsonProperty("name")
    public String functionName;

    @JsonProperty("parameters")
    public List<AbstractType> parameters;

    @JsonProperty("return_type")
    public AbstractType returnType;

    // Default constructor for Jackson
    public FunctionConfig() {}

      @JsonCreator
    public FunctionConfig(
            @JsonProperty("name") String functionName,
            @JsonProperty("parameters") List<AbstractType> parameters,
            @JsonProperty("return_type") AbstractType returnType) {
        this.functionName = functionName;
        this.parameters = parameters;
        this.returnType = returnType;
    }

    @Override
    public String toString() {
        return "FunctionConfig{" +
                "functionName='" + functionName + '\'' +
                ", parameters=" + parameters +
                ", returnType=" + returnType +
                '}';
    }
}
