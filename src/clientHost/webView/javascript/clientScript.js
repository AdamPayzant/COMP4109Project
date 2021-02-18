const socket = io();

socket.on('connect', () => {
  socket.send('Hello!');
  alert("Test");
});

socket.on('message', data => {
  console.log(data);
});