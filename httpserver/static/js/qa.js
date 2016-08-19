function updateRequestHistory() {
  blocksize = $("#datasize").val()
  var urlRequestHistory = "/api/request/history?blocksize="+blocksize;
  $.getJSON( urlRequestHistory, {
    format: "json"
  })
  .done(function(data) {
    var data_reverse = data.slice(0).reverse();
    $(".bodycontent").remove()
    $.each(data_reverse, function(i, item) {
      $("#tableData").append(
        "<tbody id=\""+item.rid+"\" class=\"bodycontent\">"+
          "<tr class=\"trcontent\">"+
            "<th class=\"text-center\">"+
              "<a href=\"javascript:void(0)\"><span class=\"glyphicon glyphicon-plus\" aria-hidden=\"true\" onclick=\"getRequestDetails("+item.rid+")\"></span></a>"+
            "</th>"+
            "<td>"+item.time+"</td>"+
            "<td>"+item.rid+"</td>"+
            "<td>"+item.supervisor+"</td>"+
            "<td>"+item.module+"</td>"+
            "<td>"+item.moduleversion+"</td>"+
            "<td>"+item.status+"</td>"+
            "<td>"+item.duration_ms+"</td>"+
          "</tr>"+
        "</tbody>")
    });
  });
}

function getRequestDetails(rid) {
  var urlRequestDetails = "/api/request/details/"+rid;
  $.getJSON(urlRequestDetails, {
    format: "json"
  })
  .done(function(data) {
    // if (data.length == 0) {
    //   alert("No details data found!")
    //   return
    // }

    var details =
      "<table class=\"table table-condensed table-bordered table-striped table-hover\">"+
      "    <thead>"+
      "      <tr>"+
      "        <th>Setp</th>"+
      "        <th>Host</th>"+
      "        <th>Supervisor</th>"+
      "        <th>Module</th>"+
      "        <th>ModuleVer</th>"+
      "        <th>Status</th>"+
      "        <th>Topic</th>"+
      "        <th>Partition</th>"+
      "        <th>Offset</th>"+
      "        <th>Duration [ms]</th>"+
      "        </tr>"+
      "    </thead>";

    $.each(data, function(i, item) {
      //console.log(item)
      console.log(details)
      details +=
        "<tbody id=\""+data[0].rid+"\" class=\"bodycontent\">"+
        "  <tr class=\"trcontent\">"+
        "    <td>"+item.step+"</td>"+
        "    <td>"+item.host+"</td>"+
        "    <td>"+item.supervisor+"</td>"+
        "    <td>"+item.module+"</td>"+
        "    <td>"+item.moduleversion+"</td>"+
        "    <td>"+item.status+"</td>"+
        "    <td>"+item.topic+"</td>"+
        "    <td>"+item.partition+"</td>"+
        "    <td>"+item.offset+"</td>"+
        "    <td>"+item.duration_ms+"</td>"+
        "  </tr>"+
        "</tbody>";
    });
    details += "</table>";

    $("#"+rid).append(
      "    <tr class=\"trcontent\">"+
      "        <th class=\"text-center\">"+
      "            <span class=\"glyphicon glyphicon-info-sign\" aria-hidden=\"true\"></span>"+
      "        </th>"+
      "        <td colspan=\"7\">"+details+"</td>"+
      "    </tr>")
  });
}
