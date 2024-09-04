package pages

const GA_CODE = `<script>
(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
})(window,document,'script','https://www.google-analytics.com/analytics.js','ga');
ga('create','UA-69818548-1','auto');
ga('send','pageview');
//
(function(l){
	var e=/(#|&)userId=(\w+)(&|$)/,
	    m=l.hash.match(e);
	if(m){
		var userId = m[2];
      ga('set','userId',userId);

      var xhReq = new XMLHttpRequest();
		xhReq.open(http.MethodGet, "/api4debtus/user?id="+userId, false);
		xhReq.send(null)

      l.replace(l.hash.replace(e,m[3]));
    }
})(location)
</script>
<script type="text/javascript">
window._urq = window._urq || [];
_urq.push(['initSite', '6ed87444-76e3-43ee-8b6e-fd28d345e79c']);
(function() {
var ur = document.createElement('script'); ur.type = 'text/javascript'; ur.async = true;
ur.src = ('https:' == document.location.protocol ? 'https://cdn.userreport.com/userreport.js' : 'http://cdn.userreport.com/userreport.js');
var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ur, s);
})();
</script>
`
