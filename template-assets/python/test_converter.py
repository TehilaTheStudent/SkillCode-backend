import unittest
import ds_utils
from ds_utils import TreeNode, Graph, ListNode
from converter import listy_to_type  # Import your `listy_to_type` function here


class TestListyToTypeIntegration(unittest.TestCase):
    def test_tree_node(self):
        stringy_input = "[1, 2, 3, 1, 1, 4, 5]"
        abstract_type = {
            "type": "TreeNode",
            "type_children": {
                "type": "Integer"
            }
        }

        # Test parsing and exporting TreeNode
        tree = listy_to_type(stringy_input, abstract_type)
        self.assertTrue(isinstance(tree, TreeNode))
        self.assertEqual(ds_utils.export_tree(tree), [1, 2, 3, 1, 1, 4, 5])

    def test_graph(self):
        stringy_input = "[[1, 2], [2, 3], [3, 1]]"
        abstract_type = {
            "type": "Graph",
            "type_children": {
                "type": "Integer"
            }
        }

        # Test parsing and exporting Graph
        graph = listy_to_type(stringy_input, abstract_type)
        self.assertTrue(isinstance(graph, Graph))
        self.assertEqual(ds_utils.export_graph(graph), [[1, 2], [2, 3], [3, 1]])

    def test_list_node(self):
        stringy_input = "[1, 2, 3, 4]"
        abstract_type = {
            "type": "ListNode",
            "type_children": {
                "type": "Integer"
            }
        }

        # Test parsing and exporting ListNode
        linked_list = listy_to_type(stringy_input, abstract_type)
        self.assertTrue(isinstance(linked_list, ListNode))
        self.assertEqual(ds_utils.export_linked_list(linked_list), [1, 2, 3, 4])

    def test_array(self):
        stringy_input = "[1, 2, 3, 4]"
        abstract_type = {
            "type": "Array",
            "type_children": {
                "type": "Integer"
            }
        }

        # Test parsing and exporting Array
        array = listy_to_type(stringy_input, abstract_type)
        self.assertTrue(isinstance(array, list))
        self.assertEqual(array, [1, 2, 3, 4])

    def test_matrix(self):
        stringy_input = "[[1, 2], [3, 4], [5, 6]]"
        abstract_type = {
            "type": "Matrix",
            "type_children": {
                "type": "Integer"
            }
        }

        # Test parsing and exporting Matrix
        matrix = listy_to_type(stringy_input, abstract_type)
        self.assertTrue(isinstance(matrix, list))
        self.assertEqual(matrix, [[1, 2], [3, 4], [5, 6]])

    # def test_nested_tree_in_graph(self):
    #   # Input: Graph of trees
    #     stringy_input = "[[[1, 2, 1], [2, 1, 1]], [[2, 1, 1], [3, 1, 1]]]"
    #     abstract_type = {
    #         "type": "Graph",
    #         "type_children": {
    #             "type": "TreeNode",
    #             "type_children": {
    #                 "type": "Integer"
    #             }
    #         }
    #     }


    #     # Test parsing and exporting nested TreeNode in Graph
    #     graph = listy_to_type(stringy_input, abstract_type)
    #     self.assertTrue(isinstance(graph, Graph))
    #     exported_graph = ds_utils.export_graph(graph)
    #     self.assertEqual(
    #         exported_graph,
    #         [[TreeNode(1), TreeNode(2)]]
    #     )

    def test_invalid_type(self):
        stringy_input = "[1, 2, 3]"
        abstract_type = {"type": "Unknown"}
        with self.assertRaises(ValueError):
            listy_to_type(stringy_input, abstract_type)


if __name__ == "__main__":
    unittest.main()
