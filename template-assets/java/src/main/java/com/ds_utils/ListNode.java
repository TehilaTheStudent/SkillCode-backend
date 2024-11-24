package com.ds_utils;

import java.util.Objects;

public class ListNode {
    public Integer val;
    public ListNode next;

    public ListNode(Integer val) {
        this.val = val;
    }

    @Override
    public boolean equals(Object obj) {
        if (!(obj instanceof ListNode))
            return false;
        ListNode current = this;
        ListNode other = (ListNode) obj;

        while (current != null && other != null) {
            if (!Objects.equals(current.val, other.val))
                return false;
            current = current.next;
            other = other.next;
        }

        return current == null && other == null;
    }

    @Override
    public int hashCode() {
        int hash = 0;
        ListNode current = this;
        while (current != null) {
            hash = Objects.hash(hash, current.val);
            current = current.next;
        }
        return hash;
    }

    @Override
    public String toString() {
        return "ListNode{" + "val=" + val + ", next=" + next + '}';
    }
}
