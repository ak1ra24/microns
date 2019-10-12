<!DOCTYPE html>
<html lang="ja">
<head>
  <meta charset="utf-8">
  <title>MicroNs Network Topology Web View</title>
  <style>
  </style>
</head>
<body style="background-color:#ddd">
  <script src="http://d3js.org/d3.v3.js" charset="utf-8"></script>
  <script src="http://marvl.infotech.monash.edu/webcola/cola.v3.min.js" charset="utf-8"></script>
  <script>
    var w = 800;
    var h = 800;

    console.log({{ .D3Nodes }})
    var nodes = {{ .D3Nodes }};
    var links = {{ .D3Links }};

    var force = cola.d3adaptor()
                .nodes(nodes)
                .links(links)
                .linkDistance(300)
				.avoidOverlaps(true)
				.symmetricDiffLinkLengths(50)
                .size([w, h])

    force.start()

    var svg = d3.select("body").append("svg").attr({width:w, height:h});
    var link = svg.selectAll(".link")
                .data(links)
                .enter()
                .append("line")
                .attr("class", "link")
                .style({stroke: "#ccc", "stroke-width": 1});
    var node = svg.selectAll(".node")
                .data(nodes)
                .enter()
                .append("g")
                .attr("class", "node")
                .call(force.drag);
    node.append("image")
        .attr("xlink:href", function(d) { 
            var image = "../images/" + d.types + ".png"
            return image
        })
        .attr("x", -10)
        .attr("y", -10)
        .attr("width", 75)
        .attr("height", 75)
    node.append("text")
        .attr({
            "text-anchor":"middle",
            "fill":"black"
        })
        .style({"font-size":11})
        .text(function(d) { return d.label; });

    force.on("tick", function() {
      link.attr({x1: function(d) { return d.source.x; },
                 y1: function(d) { return d.source.y; },
                 x2: function(d) { return d.target.x; },
                 y2: function(d) { return d.target.y; }});
      node.attr("transform", function(d) { return "translate(" + d.x + "," + d.y + ")"; });
    });

  </script>
</body>
</html>
