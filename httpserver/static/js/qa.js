(function () {
  $(window).load(function() {
    updateRequestHistory();
  });
})();

function updateRequestHistory() {
  var status_green = "#88CF85"
  var status_red = "#F87070"
  var duration_green = "#48B444"
  var duration_orange = "#F4B27B"
  var duration_red = "#F53636"
  var blocksize = $("#datasize").val()
  var urlRequestHistory = "/api/request/history?blocksize="+blocksize;
  $.getJSON( urlRequestHistory, {
    format: "json"
  })
  .done(function(data) {
    var count = data.length
    var data_reverse = data.slice(0).reverse();
    $(".bodycontent").remove()
    $.each(data_reverse, function(i, item) {
      var status_color = item.status == 200 ? status_green : status_red;
      var duration_color;
      if (item.duration_ms < 1000) {
        duration_color = duration_green;
      } else if (item.duration_ms >= 1000 && item.duration_ms < 1500) {
        duration_color = duration_orange;
      } else {
        duration_color = duration_red;
      }
      $("#tableData").append(
        "<tbody id=\""+item.rid+"\" class=\"bodycontent\">"+
          "<tr class=\"trcontent\">"+
            "<th class=\"text-center\">"+
              "<a href=\"javascript:void(0)\"><span id=\"icon-"+item.rid+"\" class=\"glyphicon glyphicon-plus\" aria-hidden=\"true\" onclick=\"getRequestDetails("+item.rid+")\"></span></a>"+
            "</th>"+
            "<td>"+item.time+"</td>"+
            "<td>"+item.rid+"</td>"+
            "<td>"+item.supervisor+"</td>"+
            "<td>"+item.module+"</td>"+
            "<td>"+item.moduleversion+"</td>"+
            "<td><span style=\"color:"+duration_color+"\">"+item.duration_ms+"</span></td>"+
            "<td bgcolor=\""+status_color+"\">"+item.status+"</td>"+
          "</tr>"+
        "</tbody>");
    });
    $("#tableData").append(
      "<tbody class=\"bodycontent\">"+
      "  <tr class=\"trcontent\">"+
      "    <td colspan=\"8\" class=\"text-center\"><h6>Found "+count+" records</h6></td>"+
      "  </tr>"+
      "</tbody>");
  });
}

function getRequestDetails(rid) {

  if ($("#tbase-"+rid).length) {
    $("#icon-"+rid).removeClass("glyphicon-minus").addClass("glyphicon-plus");
    $("#tbase-"+rid).remove();
    return;
  }

  var urlRequestDetails = "/api/request/details/"+rid;
  $.getJSON(urlRequestDetails, {
    format: "json"
  })
  .done(function(data) {
    var details =
      "<table class=\"table table-condensed table-bordered table-striped table-hover\">"+
      "    <thead>"+
      "      <tr>"+
      "        <th>Step</th>"+
      "        <th>RequestId</th>"+
      "        <th>Host</th>"+
      "        <th>Supervisor</th>"+
      "        <th>Module</th>"+
      "        <th>ModuleVer</th>"+
      "        <th>Message</th>"+
      "        <th>Topic</th>"+
      "        <th>Partition</th>"+
      "        <th>Offset</th>"+
      "        <th>Duration [ms]</th>"+
      "        <th>Status</th>"+
      "        </tr>"+
      "    </thead>";

    $.each(data, function(i, item) {
      details +=
        "<tbody id=\"tb-"+item.rid+"\" class=\"bodycontent\">"+
        "  <tr class=\"trcontent\">"+
        "    <td>"+item.step+"</td>"+
        "    <td>"+item.rid+"</td>"+
        "    <td>"+item.host+"</td>"+
        "    <td>"+item.supervisor+"</td>"+
        "    <td>"+item.module+"</td>"+
        "    <td>"+item.moduleversion+"</td>"+
        "    <td>"+item.message+"</td>"+
        "    <td>"+item.topic+"</td>"+
        "    <td>"+item.partition+"</td>"+
        "    <td>"+item.offset+"</td>"+
        "    <td>"+item.duration_ms+"</td>"+
        "    <td>"+item.status+"</td>"+
        "  </tr>"+
        "</tbody>";
    });
    details += "</table>";
    $("#icon-"+rid).removeClass("glyphicon-plus").addClass("glyphicon-minus");
    $("#"+rid).append(
      "    <tr id=\"tbase-"+rid+"\" class=\"trcontent\">"+
      "        <th class=\"text-center\">"+
      "        </th>"+
      "        <td colspan=\"8\">"+details+"</td>"+
      "    </tr>")
  });
}
