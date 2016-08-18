function updateRequestHistory() {
  blocksize = $("#datasize").val()
  var urlRequestHistory = "/api/request/history?blocksize="+blocksize;
  $.getJSON( urlRequestHistory, {
    format: "json"
  })
  .done(function( data ) {
    $(".bodycontent").remove()
    $.each( data, function( i, item ) {
      console.log(item)
      $("#tableData").append(
        "<tbody class=\"bodycontent\">"+
          "<tr class=\"trcontent\">"+
            "<th class=\"text-center\">"+
              "<a href=\"#\"><span class=\"glyphicon glyphicon-plus\" aria-hidden=\"true\" onclick=\"getRequestDetails("+item.rid+")\"></span></a>"+
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
  alert("details: "+rid)
}
