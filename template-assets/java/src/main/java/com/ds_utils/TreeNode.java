package com.ds_utils;
import java.util.Objects;

public class TreeNode {
    public Integer val;
    public TreeNode left;
    public TreeNode right;

    public TreeNode(Integer val) {
        this.val = val;
    }

    @Override
    public boolean equals(Object obj) {
        if (!(obj instanceof TreeNode)) return false;
        TreeNode other = (TreeNode) obj;
        return Objects.equals(val, other.val)
                && Objects.equals(left, other.left)
                && Objects.equals(right, other.right);
    }

    @Override
    public int hashCode() {
        return Objects.hash(val, left, right);
    }

    @Override
    public String toString() {
        return "TreeNode{" + "val=" + val + ", left=" + left + ", right=" + right + '}';
    }
}
