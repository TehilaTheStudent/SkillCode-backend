package com.evaluation;

import java.util.List;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

public class TestCase {
    @JsonProperty("parameters")
    public List<String> parameters;

    @JsonProperty("expected_output")
    public String expectedOutput;

    // Default constructor for Jackson
    public TestCase() {}

       @JsonCreator
    public TestCase(
            @JsonProperty("parameters") List<String> parameters,
            @JsonProperty("expected_output") String expectedOutput) {
        this.parameters = parameters;
        this.expectedOutput = expectedOutput;
    }

    @Override
    public String toString() {
        return "TestCase{" +
                "parameters=" + parameters +
                ", expectedOutput='" + expectedOutput + '\'' +
                '}';
    }
}
