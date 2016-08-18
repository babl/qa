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
    var details
    $.each(data, function(i, item) {
      console.log(item)
      details += ''
    });


    $("#"+rid).append(
      "    <tr class=\"trcontent\">"+
      "        <th class=\"text-center\">"+
      "            <span class=\"glyphicon glyphicon-info-sign\" aria-hidden=\"true\"></span>"+
      "        </th>"+
      "        <td colspan=\"7\"></td>"+
      "    </tr>")
  });
}
