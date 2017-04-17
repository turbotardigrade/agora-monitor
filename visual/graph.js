const svg = d3.select('svg');
const radius = 20;

function createNode(id, x, y, col) {
  return { id: id, radius, fill: col, x, y };
}

function createLink(source, target, stroke) {
  return { source, target, stroke };
}

const data = {
  nodes:[
    createNode('a', 50, 100, 'red'),
    createNode('b', 200, 100, 'green'),
    createNode('c', 350, 100, 'green'),
    createNode('d', 500, 200, 'green'),
    createNode('e', 650, 200, 'green'),
    createNode('f', 50, 300, 'green'),
    createNode('g', 200, 300, 'green'),
    createNode('h', 350, 300, 'green')
  ],
  links: [
    createLink('a', 'b', 'red'),
    createLink('a', 'g', 'red'),
    createLink('a', 'h', 'red'),
    createLink('f', 'g', 'green'),
    createLink('g', 'h', 'green'),
    createLink('g', 'c', 'green'),
    createLink('b', 'c', 'green'),
    createLink('c', 'h', 'green'),
    createLink('c', 'd', 'green'),
    createLink('h', 'd', 'green'),
    createLink('d', 'e', 'green')
  ]
};

const gLinks = svg.append('g').attr('class', 'links');
const gNodes = svg.append('g').attr('class', 'nodes');


// Takes an object representing the state of the network and draws it
// Warning, possible mindfuck
// Helpful links:
// https://bost.ocks.org/mike/selection/
// https://bost.ocks.org/mike/join/
// https://bl.ocks.org/mbostock/3808218
function draw(data) {
  // DATA JOIN
  // Join new data with old elements, if any.

  const link = gLinks
    .selectAll('line')
    .data(data.links);
  const node = gNodes
    .selectAll('g')
    .data(data.nodes);

  // ENTER
  // Create new elements as needed.
  //
  // ENTER + UPDATE
  // After merging the entered elements with the update selection,
  // apply operations to both.

  const nodeEnter = node.enter().append('g');

  nodeEnter.append('circle');
  nodeEnter.append('text');
  nodeEnter.attr('transform', d => `translate(${d.x},${d.y})`);

  const nodeUpdate = nodeEnter.merge(node);

  nodeUpdate.select('circle')
	    .attr('r', d => d.radius)
	    .style('fill', d => d.fill);
  nodeUpdate.select('text')
	    .attr('text-anchor', 'middle')
	    .attr('dy', '0.35em') // http://stackoverflow.com/a/8684888/4131237
	    .text(d => d.id);

  const linkEnter = link.enter().append('line');

  // HACK
  // I had to use <g> elements as a container for <circle> and <text>
  // But <g> elements don't have coordinate attributes
  // So I tried getting their rendered coordinates using getBoundingClientRect()
  // However even that doesn't have x and y attributes
  // but they have left and top which seem to be the same thing according to
  // https://developer.mozilla.org/en-US/docs/Mozilla/Tech/XPCOM/Reference/Interface/nsIDOMClientRect
  // This is probably why the link drawing is a bit off
  // TODO figure out a better solution
  linkEnter.attr('x1', d => gNodes.selectAll('g')
				  .filter(g => g.id === d.source)
				  .node()
				  .getBoundingClientRect().left+0.6*radius);
  linkEnter.attr('y1', d => gNodes.selectAll('g')
				  .filter(g => g.id === d.source)
				  .node()
				  .getBoundingClientRect().top+0.6*radius);
  linkEnter.attr('x2', d => gNodes.selectAll('g')
				  .filter(g => g.id === d.target)
				  .node()
				  .getBoundingClientRect().left+0.6*radius);
  linkEnter.attr('y2', d => gNodes.selectAll('g')
				  .filter(g => g.id === d.target)
				  .node()
				  .getBoundingClientRect().top+0.6*radius);

  const linkUpdate = linkEnter.merge(link);
  linkUpdate.style('stroke', d => d.stroke);

  // EXIT
  // Remove old elements as needed.
  node.exit().remove();
  link.exit().remove();
};

draw(data);

function getNode(id) {
  for (i=0; i<data.nodes.length; ++i) {
    if (data.nodes[i].id === id) {
      return data.nodes[i]
    }
  }

  console.log("Error: ID "+id+" not found");
}

function getLink(id1, id2) {
  for (i=0; i<data.links.length; ++i) {
    var link = data.links[i];
    var isLink1 = link.source === id1 && link.target === id2;
    var isLink2 = link.source === id2 && link.target === id1;
    if (isLink1 || isLink2) {
      return link;
    }
  }

  conosle.log("Error: Link between "+id1+" and "+id2+" not found");
}

// Example of how to redraw with new data
// Just call draw() again with new data
setTimeout(function() {
  getNode('b').fill = 'black';
  getLink('a', 'b').stroke = 'black';

  draw(data);
}, 2000);
