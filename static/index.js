var app = new Vue({ 
    el: '#app',
    data: {
        message: 'Hello Vue!'
    }
});


var app2 = new Vue({
    el: '#app-2',
    data: {
      message: 'Hai caricato questa pagina il ' + new Date().toLocaleString()
    }
  })

var appdomande = new Vue({
    el: '#domande',
    data: {
      visibile: true
    },
    methods:{
      callFunction: function () {
          var v = this;
          setTimeout(function () {
              v.visibile = false;
          }, 900000);
      }
  },
  mounted () {
    this.callFunction()
  }
});
 
var app3 = new Vue({
  el: '#app-3',
  data: {
    seen: true
  }
})

var app4 = new Vue({
  el: '#app-4',
  data: {
    todos: [
      { text: 'Learn JavaScript' },
      { text: 'Learn Vue' },
      { text: 'Build something awesome' }
    ]
  }
})