const svg = d3.select('svg');
const radius = 20;
var link = null;
var node = null;

var nodes_data = [];
var links_data = [];

var nodes = null;
$.get( "/nodes", function( data ) {
  nodes = data;
  init();
});

function init() {
  // Create nodes
  for (var i=0; i<nodes.length; ++i) {
    var elem = nodes[i];

    if (elem.is_spammer) {
      nodes_data.push(createNode(elem.id, 'rgb(255, 0, 0)'))
    } else {
      nodes_data.push(createNode(elem.id, 'rgb(0, 255, 0)'))
    };
  }

  // Connect all nodes together
  for (var i=0; i<nodes_data.length-1; ++i) {
    for (var j=i+1; j<nodes_data.length; ++j) {
      var s = nodes_data[i].id;
      var t = nodes_data[j].id;

      var col = 'green';
      if (nodes[i].is_spammer || nodes[j].is_spammer) {
	col = 'red';
      }
      
      links_data.push(createLink(s, t, col));
    }
  }

  link = svg.selectAll(".link").data(links_data).enter().append("line");
  node = svg.selectAll(".node").data(nodes_data).enter().append("g");

  draw();
  setTimeout(update, 3000);
}

// @TODO spammer should have red border

function update() {
  console.log("update");

  $.get( "/monitor", function(data) {
    Object.keys(data).forEach(function(k) {
      var n = getNode(k);
      var d = data[k];
      
      if (!d.healthy) {
	 n.fill = 'gray';
	 colorLinksOfNode(n.id, 'gray');
	 return
      }

      var red = Math.min(Math.round(d.spam_ratio*5*255), 255);
      var green = Math.min(Math.round((1-d.spam_ratio*5)*255), 255);
      n.fill = "rgb("+red+","+green+",0)";

      // For testing
      /* if (k != "QmaXjNr6ZER1vpCZ7LPrsABCwtMRbwWfDVRHPcwJmfKwSs" && k != "QmPxeRd41wZYg8DuDxbV6mCfHRRTFJua5qwihmWHahmUQ1") {
	 d.blacklist = [{"QmPxeRd41wZYg8DuDxbV6mCfHRRTFJua5qwihmWHahmUQ1":20.0}, {"QmaXjNr6ZER1vpCZ7LPrsABCwtMRbwWfDVRHPcwJmfKwSs":20.0}];
	 } */
      
      if (d.blacklist) {
	d.blacklist.forEach(function(elem){
	  Object.keys(elem).forEach(function(id){
	    removeLink(k, id);
	  });
	});
      }
      
    });
    
    draw();
  });

  // Example of how to redraw with new data
  // Just call draw() again with new data
  setTimeout(update, 2000);
}

// Takes an object representing the state of the network and draws it
// Warning, possible mindfuck
// Helpful links:
// https://bost.ocks.org/mike/selection/
// https://bost.ocks.org/mike/join/
// https://bl.ocks.org/mbostock/3808218
function draw() {
  var width = 960;
  var height = 500;

  var force = d3.layout.force()
		.gravity(.05)
		.distance(200)
		.charge(-400)
		.size([width, height]);

  force.nodes(nodes_data).links(links_data).start();

  link.attr("class", "link")
      .style('stroke', d => d.stroke)
      .style("stroke-width", function(d) { return Math.sqrt(d.weight); });

  node.attr("class", "node").call(force.drag);

  node.append("circle")
      .attr('r', d => d.radius)
      .style('fill', d => d.fill)
      .style('stroke', d => d.stroke);

  node.append("text")
      .attr('text-anchor', 'middle')
      .attr('dy', '0.35em') // http://stackoverflow.com/a/8684888/4131237
      .text(d => d.id.substr(d.id.length-3));

  force.on("tick", function() {
    link.attr("x1", function(d) { return d.source.x; })
        .attr("y1", function(d) { return d.source.y; })
        .attr("x2", function(d) { return d.target.x; })
        .attr("y2", function(d) { return d.target.y; });

    node.attr("transform", function(d) { return "translate(" + d.x + "," + d.y + ")"; });
  });
  
  svg.selectAll('.link').data(links_data).exit().remove();
  node.selectAll('g').data(nodes_data).exit().remove();
}

//////////////////////////////////////////////////////////////////////
// Graph helpers

function getNode(id) {
  for (var i=0; i<nodes_data.length; ++i) {
    if (nodes_data[i].id === id) {
      return nodes_data[i]
    }
  }

  console.log("Error: ID "+id+" not found");
}

function getLink(id1, id2) {
  for (var i=0; i<links_data.length; ++i) {
    var link = links_data[i];
    var isLink1 = link.source.id == id1 && link.target.id == id2;
    var isLink2 = link.source.id == id2 && link.target.id == id1;
    if (isLink1 || isLink2) {
      return link;
    }
  }

  console.log("Error: Link between "+id1+" and "+id2+" not found");
}

function removeLink(id1, id2) {
  for (var i=0; i<links_data.length; ++i) {
    var link = links_data[i];
    var isLink1 = link.source.id == id1 && link.target.id == id2;
    var isLink2 = link.source.id == id2 && link.target.id == id1;
    if (isLink1 || isLink2) {
      // removes the link by 'drawing' a link to the same node
      links_data[i].source = links_data[i].target;
      return
    }
  }

  // console.log("Error: Link between "+id1+" and "+id2+" not found");
}

function createNode(id, col) {
  return { id: id, radius, fill: col, stroke: col};
}

// Connect all nodes together
function findNodeIndex(id) {
  for (var i=0; i<nodes_data.length; i++) {
    if (nodes_data[i].id === id) {
      return i;
    }
  }

  console.log("Error: ID "+id+" not found");
  return -1;
}

function createLink(source, target, col) {
  var s = findNodeIndex(source);
  var t = findNodeIndex(target);
  return { source: s, target: t, weight: 3};
}

function colorLinksOfNode(nodeID, col) {
  var links = []
  for (var i=0; i<links_data.length; ++i) {
    var isSource = links_data[i].source.id == nodeID
    var isTarget = links_data[i].target.id == nodeID
    if (isSource || isTarget) {
      links_data[i].stroke = col;
    }
  }
}
