const svg = d3.select('svg')

  const data = {
    nodes: [{
      id: 'a',
      radius: 20,
      fill: 'red',
      x: 50,
      y: 100
    }, {
      id: 'b',
      radius: 20,
      fill: 'green',
      x: 200,
      y: 100
    }, {
      id: 'c',
      radius: 20,
      fill: 'green',
      x: 350,
      y: 100
    }, {
      id: 'd',
      radius: 20,
      fill: 'green',
      x: 500,
      y: 200
    }, {
      id: 'e',
      radius: 20,
      fill: 'green',
      x: 650,
      y: 200
    }, {
      id: 'f',
      radius: 20,
      fill: 'green',
      x: 50,
      y: 300
    }, {
      id: 'g',
      radius: 20,
      fill: 'green',
      x: 200,
      y: 300
    }, {
      id: 'h',
      radius: 20,
      fill: 'green',
      x: 350,
      y: 300
    }],
    links: [{
      source: 'a',
      target: 'b',
      stroke: 'red'
    }, {
      source: 'a',
      target: 'g',
      stroke: 'red'
    }, {
      source: 'a',
      target: 'h',
      stroke: 'red'
    }, {
      source: 'f',
      target: 'g',
      stroke: 'green'
    }, {
      source: 'g',
      target: 'h',
      stroke: 'green'
    }, {
      source: 'g',
      target: 'c',
      stroke: 'green'
    }, {
      source: 'b',
      target: 'c',
      stroke: 'green'
    }, {
      source: 'c',
      target: 'h',
      stroke: 'green'
    }, {
      source: 'c',
      target: 'd',
      stroke: 'green'
    }, {
      source: 'h',
      target: 'd',
      stroke: 'green'
    }, {
      source: 'd',
      target: 'e',
      stroke: 'green'
    }]
  }

const gNodes = svg.append('g')
		  .attr('class', 'nodes')
  const gLinks = svg.append('g')
		    .attr('class', 'links')

  // Takes an object representing the state of the network and draws it
  // Warning, possible mindfuck
  // Helpful links:
  // https://bost.ocks.org/mike/selection/
  // https://bost.ocks.org/mike/join/
  // https://bl.ocks.org/mbostock/3808218
  function draw(data) {
    // DATA JOIN
    // Join new data with old elements, if any.

    const node = gNodes
    .selectAll('g')
    .data(data.nodes)
      const link = gLinks
    .selectAll('line')
    .data(data.links)

      // ENTER
      // Create new elements as needed.
      //
      // ENTER + UPDATE
      // After merging the entered elements with the update selection,
      // apply operations to both.

      const nodeEnter = node.enter().append('g')

      nodeEnter
			    .append('circle')
      nodeEnter
			    .append('text')
      nodeEnter
			    .attr('transform', d => `translate(${d.x},${d.y})`)

      const nodeUpdate = nodeEnter.merge(node)   

      nodeUpdate.select('circle')
				  .attr('r', d => d.radius)
				  .style('fill', d => d.fill)
      nodeUpdate.select('text')
				  .attr('text-anchor', 'middle')
				  .attr('dy', '0.35em') // http://stackoverflow.com/a/8684888/4131237
				  .text(d => d.id)

      const linkEnter = link.enter().append('line')

      // HACK
      // I had to use <g> elements as a container for <circle> and <text>
      // But <g> elements don't have coordinate attributes
      // So I tried getting their rendered coordinates using getBoundingClientRect()
      // However even that doesn't have x and y attributes
      // but they have left and top which seem to be the same thing according to
      // https://developer.mozilla.org/en-US/docs/Mozilla/Tech/XPCOM/Reference/Interface/nsIDOMClientRect
      // This is probably why the link drawing is a bit off
      // TODO figure out a better solution
      linkEnter
			    .attr('x1', d => gNodes.selectAll('g')
						   .filter(g => g.id === d.source)
						   .node()
						   .getBoundingClientRect().left)
			    .attr('y1', d => gNodes.selectAll('g')
						   .filter(g => g.id === d.source)
						   .node()
						   .getBoundingClientRect().top)
			    .attr('x2', d => gNodes.selectAll('g')
						   .filter(g => g.id === d.target)
						   .node()
						   .getBoundingClientRect().left)
			    .attr('y2', d => gNodes.selectAll('g')
						   .filter(g => g.id === d.target)
						   .node()
						   .getBoundingClientRect().top)

      const linkUpdate = linkEnter.merge(link)

      linkUpdate
				  .style('stroke', d => d.stroke)

      // EXIT
      // Remove old elements as needed.
      node.exit().remove()
      link.exit().remove()
  }

draw(data)

  // Example of how to redraw with new data
  // Just call draw() again with new data
  setTimeout(function() {
    data.nodes[0].fill = 'blue'
    data.links[2].stroke = 'purple'
    
    draw(data)
  }, 2000)
