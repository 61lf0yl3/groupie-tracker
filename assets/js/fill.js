document.getElementById('searchText').addEventListener('input', function () {
  var inp = document.getElementById('searchText').value.split(' -> ');
  if (inp[1] == "artist/band") {
    document.getElementById('searchText').value = inp[0];
    document.getElementById('option1').click()
  };
  if (inp[1] == "member") {
    document.getElementById('searchText').value = inp[0];
    document.getElementById('option2').click()
  };
  if (inp[1] == "first album date") {
      document.getElementById('searchText').value = inp[0];
      document.getElementById('option3').click()
  };
  if (inp[1] == "creation date") {
      document.getElementById('searchText').value = inp[0];
      document.getElementById('option4').click()
  };
  if (inp[1] == "concert location") {
    document.getElementById('searchText').value = inp[0];
    document.getElementById('option5').click()
  };
  if (inp[1] == "concert date") {
    document.getElementById('searchText').value = inp[0];
    document.getElementById('option6').click()
};
});
