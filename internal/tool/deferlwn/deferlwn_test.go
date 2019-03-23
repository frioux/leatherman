package deferlwn

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeferLink(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Couldn't create dir: %s", err)
	}
	defer os.RemoveAll(d)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN"
        "http://www.w3.org/TR/html4/loose.dtd">
        <html>
        <head><title>Subscription required [LWN.net]</title>
        <meta name="viewport" content="width=device-width, initial-scale=1">
<meta HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=utf-8">
        <link rel="icon" href="/images/favicon.png" type="image/png">
        
        <link rel="stylesheet" href="/CSS/lwn">
<link rel="stylesheet" href="/CSS/nosub">
<link rel="stylesheet" href="/CSS/pure-min">
           <!--[if lte IE 8]>
             <link rel="stylesheet" href="/CSS/grids-responsive-old-ie-min">
           <![endif]-->
           <!--[if gt IE 8]><!-->
             <link rel="stylesheet" href="/CSS/grids-responsive-min">
           <!--<![endif]-->
           <link rel="stylesheet" href="/CSS/pure-lwn">
           
        
        </head>
        <body bgcolor="#ffffff" link="Blue" VLINK="Green" alink="Green">
        <!-- no tracking img --><a name="t"></a>
<div id="menu"><a href="/"><img src="https://static.lwn.net/images/logo/barepenguin-70.png" class="logo"
                 border="0" alt="LWN.net Logo">
           <font class="logo">LWN<br>.net</font>
           <font class="logobl">News from the source</font></a>
           <a href="/"><img src="https://static.lwn.net/images/lcorner-ss.png" class="sslogo"
                 border="0" alt="LWN"></a><div class="navmenu-container">
           <ul class="navmenu">
        <li><a class="navmenu" href="#t"><b>Content</b></a><ul><li><a href="/current/">Weekly Edition</a></li><li><a href="/Archives/">Archives</a></li><li><a href="/Search/">Search</a></li><li><a href="/Kernel/">Kernel</a></li><li><a href="/Security/">Security</a></li><li><a href="/Distributions/">Distributions</a></li><li><a href="/Calendar/">Events calendar</a></li><li><a href="/Comments/unread">Unread comments</a></li><li><hr></li><li><a href="/op/FAQ.lwn">LWN FAQ</a></li><li><a href="/op/AuthorGuide.lwn">Write for us</a></li></ul></li>
</ul></div>
</div> <!-- menu -->
<div class="pure-g not-handset" style="margin-left: 10.5em">
           <div class="not-print">
             <!-- no ads -->
           </div>
           </div>
        <div class="topnav-container">
<div class="not-handset"><form action="https://lwn.net/Login/" method="post" name="loginform"
                 class="loginform">
        <b>User:</b> <input type="text" name="Username" value="" size="8" /> <b>Password:</b> <input type="password" name="Password" size="8" /> <input type="hidden" name="target" value="/Articles/783673/" /> <input type="submit" name="submit" value="Log in" /></form> |
           <form action="https://lwn.net/subscribe/" method="post" class="loginform">
           <input type="submit" name="submit" value="Subscribe" />
           </form> |
           <form action="https://lwn.net/Login/newaccount" method="post" class="loginform">
           <input type="submit" name="submit" value="Register" />
           </form>
        </div>
               <div class="handset-only">
               <a href="/subscribe/"><b>Subscribe</b></a> /
               <a href="/Login/"><b>Log in</b></a> /
               <a href="/Login/newaccount"><b>New account</b></a>
               </div>
               </div><div class="pure-grid maincolumn">
<div class="lwn-u-1 pure-u-md-19-24">
<div class="PageHeadline">
<h1>Subscription required</h1>
</div>
<div class="ArticleText">
 <!-- offer -->
The page you have tried to view  (<a href="/Articles/783673/">The congestion-notification conflict</a>)  is currently available to LWN 
subscribers only.   <!-- exp --> Reader subscriptions are a necessary way
to fund the continued existence of LWN and the quality of its content.
<p>
                   If you are already an LWN.net subscriber, please log in
                   with the form below to read this content.
                   <p>
                   <blockquote>
                   <form action="https://lwn.net/Login/" method="post" name="loginform">
           <table class="Form">
            <tr><td><b>Username</b></td><td><input type="text" name="Username" value="" size="20" /></td></tr>
 <tr><td><b>Password</b></td><td><input type="password" name="Password" size="20" /></td></tr>
 <tr><td><b><input type="hidden" name="target" value="https://lwn.net/Articles/783673/" /></b></td><td><input type="submit" name="submit" value="Log in" /></td></tr>
</table></form>
          <script type="text/javascript">
          <!--
          document.loginform.Username.focus();
          // -->
          </script>
          
                   </blockquote>
                    <!-- login -->
<p>
Please consider <a href="/subscribe">subscribing to LWN</a>.  An LWN
subscription provides numerous benefits, including access to restricted
content and the warm feeling of knowing that you are helping to keep LWN
alive.  
<p>
(Alternatively, this item will become freely
                     available on April 4, 2019)
                     

</div> <!-- ArticleText -->
</div>
<div class="lwn-u-1 pure-u-md-1-6 not-print">

</div>
</div> <!-- pure-grid -->

        <br clear="all">
        <center>
        <P>
        <font size="-2">
        Copyright &copy; 2019, Eklektix, Inc.<BR>
        
        Comments and public postings are copyrighted by their creators.<br>
        Linux  is a registered trademark of Linus Torvalds<br>
        </font>
        </center>
        
            <script type="text/javascript">
            var gaJsHost = (("https:" == document.location.protocol) ? "https://ssl." : "http://www.");
            </script>
            <script type="text/javascript">
            try {
            var pageTracker = _gat._getTracker("UA-2039382-1");
            pageTracker._trackPageview();
            } catch(err) {}</script>
            
        </body></html>
`)
	}))
	defer ts.Close()

	buf := &bytes.Buffer{}
	bufErr := &bytes.Buffer{}
	err = run(d, strings.NewReader("[title]("+ts.URL+")"), buf, bufErr)
	assert.NoError(t, err)
	assert.Empty(t, buf.Bytes())
	assert.Empty(t, bufErr.Bytes())

	dh, err := os.Open(d)
	if !assert.NoError(t, err) {
		return
	}

	n, err := dh.Readdirnames(10)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, []string{"2019-04-05-lwn.md"}, n)

	f, err := os.Open(filepath.Join(d, "2019-04-05-lwn.md"))
	if !assert.NoError(t, err) {
		return
	}

	buf = &bytes.Buffer{}
	_, err = io.Copy(buf, f)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "[title]("+ts.URL+")\n", buf.String())
}
