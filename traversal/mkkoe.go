package traversal

import ON "github.com/fbaube/orderednodes"
      
// A sort of binary search tree where every node matters.
// 
//          30
//         /  \
//        /    \
//      10     50
//     / \    /  \
//   00  20  40  60

var Node00to60 *ON.SNord

func init() {
     var nd10 = ON.NewSNord("10")
     var nd50 = ON.NewSNord("50")
     nd10.AddKids([]ON.Norder { ON.NewSNord("10"), ON.NewSNord("10") })
     nd50.AddKids([]ON.Norder { ON.NewSNord("40"), ON.NewSNord("60") })
     var nd30 = ON.NewSNord("30")
     nd30.AddKids([]ON.Norder { nd10, nd50 })
     }

// Some HTML, where there is implicit hierarchy.
//
//   html
//     head
//       title
//         txt: The Title 
//       meta
//     body
//       h1
//         txt: The Headline 
//       p #1
//         txt: This intro para is the short desc.
//       ul
//         li
//           txt: Todo #1
//         li
//           txt: Todo #2 
//       p #2
//         txt: The outro para. 
//         

var NodeHtml *ON.SNord

func init() {
     var nodeTitl = ON.NewSNord("<title>").AddKid(ON.NewSNord("txt: The Title"))
     var nodeHead = ON.NewSNord("<head>").AddKids([]ON.Norder {
     	 	    nodeTitl, ON.NewSNord("<meta>") })
     var nodeHedg = ON.NewSNord("<h1>").AddKid(ON.NewSNord("txt: The Heading"))
     var nodePar1 = ON.NewSNord("<p>" ).AddKid(ON.NewSNord(
     	 	    "txt: This intro para is the short desc."))
     var nodeLi11 = ON.NewSNord("<li>").AddKid(ON.NewSNord("txt: Todo #1"))
     var nodeLi22 = ON.NewSNord("<li>").AddKid(ON.NewSNord("txt: Todo #2"))
     var nodeUlst = ON.NewSNord("<ul>").AddKids([]ON.Norder { nodeLi11, nodeLi22 })
     var nodePar2 = ON.NewSNord("<p>" ).AddKid(ON.NewSNord(
     	 	    "txt: The outro para."))
     var nodeBody = ON.NewSNord("<body>").AddKids([]ON.Norder {
     	 	    nodeHedg, nodePar1, nodeUlst, nodePar2 })
     NodeHtml = ON.NewSNord("<html>")
     NodeHtml.AddKids([]ON.Norder { nodeHead, nodeBody })
     }