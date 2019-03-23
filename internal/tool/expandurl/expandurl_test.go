package expandurl

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
<!DOCTYPE html>
<html class="no-js" lang="en-US">
<head>
<meta name="generator" content="Hugo 0.49" />
<title>fREW Schmidt&#39;s Foolish Manifesto</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta content="text/html; charset=UTF-8" http-equiv="Content-Type">
<link href="/index.xml" rel="alternate" type="application/rss+xml" title="fREW Schmidt&#39;s Foolish Manifesto" />
<link href="/static/css/bootstrap.min.css" rel="stylesheet" />
<link href="/static/css/styles.css" rel="stylesheet" />
<script type="568c67a7cf35485aee2a5b8b-text/javascript" src="/static/js/foolish.js"></script>
<link href="/static/css/fonts.css" rel="stylesheet" />
</head>
<body>
<nav class="navbar navbar-inverse navbar-fixed-top" role="navigation">
<div class="navbar-header">
<button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-ex1-collapse">
<span class="sr-only">Toggle navigation</span>
<span class="icon-bar"></span>
<span class="icon-bar"></span>
<span class="icon-bar"></span>
</button>
<a class="navbar-brand" href="/">fREW Schmidt&#39;s Foolish Manifesto</a>
</div>
<div class="collapse navbar-collapse navbar-ex1-collapse">
<ul class="nav navbar-nav">
</ul>
<ul class="nav navbar-nav navbar-right">
<li><a class="sigil" href="https://github.com/frioux">♑</a></li>
<li><a class="sigil" href="https://twitter.com/frioux">♍</a></li>
<li><a class="sigil" href="https://secure.flickr.com/photos/frew">♌</a></li>
<li><a class="sigil" href="https://www.linkedin.com/in/frew-schmidt-567216120">🔗</a></li>
<li><a class="sigil" id="toggleTheme" onClick="if (!window.__cfRLUnblockHandlers) return false; toggleTheme()" href="#" data-cf-modified-568c67a7cf35485aee2a5b8b-="">💡</a></li>
</ul>
</div>
</nav>
<div class="container" id="main">
<div class="row">
<div class="col-md-12">
<section id="main">
<div>
<h1><a href="/posts/sorting-books/">Sorting Books</a></h1>
<p>I wrote a little program to sort lists of books.</p>
<p></p>
<br /><br />
Posted Thu, Mar 21, 2019
<h1><a href="/posts/automating-email/">Automating Email</a></h1>
<p>I just automated a couple common email tasks.</p>
<p></p>
<br /><br />
Posted Mon, Mar 18, 2019
<h1><a href="/posts/how-to-add-a-subscription-mode-to-your-blog/">How to Add a Subscription Service to Your Blog</a></h1>
<p>I used to use a service to email subscribers updates to my blog. The service
broke, but I automated my way around it.</p>
<p></p>
<br /><br />
Posted Thu, Mar 7, 2019
<h1><a href="/posts/fixing-buggy-haskell-programs-with-golang/">Fixing Buggy Haskell Programs with Go</a></h1>
<p>I recently ran into a stupid bug in a program written in Haskell and found it
much easier to paper over with a few lines of Go than to properly fix.</p>
<p></p>
<br /><br />
Posted Wed, Feb 27, 2019
<h1><a href="/posts/learning-day-2-diy-games/">Learning Day 2: DIY Games</a></h1>
<p>Today I did my second Learning Day; the subject was DIY Games.</p>
<p></p>
<br /><br />
Posted Sat, Feb 23, 2019
<h1><a href="/posts/busting-cloudflare-cache/">Busting the Cloudflare Cache</a></h1>
<p>I automated blowing the cache for this blog. Read on to see how I did it.</p>
<p></p>
<br /><br />
Posted Wed, Feb 20, 2019
 <h1><a href="/posts/graphviz/">graphviz describing multi-stage docker builds</a></h1>
<p>I recently decided I should learn to use Graphviz more, as a great tool for
making certain kinds of plots. Less than a week later a great use case
surfaced.</p>
<p></p>
<br /><br />
Posted Mon, Feb 11, 2019
<h1><a href="/posts/amygdala/">Amygdala</a></h1>
<p>This past weekend I started re-creating a tool I used to have, using new tools,
techniques, and infrastructure. The tool allows, at least, adding to my own
todo list via SMS. It&rsquo;s working great!</p>
<p></p>
<br /><br />
Posted Tue, Feb 5, 2019
<h1><a href="/posts/deploying-to-kubernetes-at-ziprecruiter/">Deploying to Kubernetes at ZipRecruiter</a></h1>
<p>At <a href="https://www.ziprecruiter.com/hiring/technology">ZR</a> we are working hard to
get stuff migrated to Kubernetes, and a big part of that is our cicd pipeline.
We have that stable enough that I can explain the major parts.</p>
<p></p>
<br /><br />
Posted Wed, Jan 30, 2019
<h1><a href="/posts/full-text-search-for-ebooks/">Full Text Search for ebooks</a></h1>
<p>This past weekend I did a learning day that inspired me to try SQLite for
indexing my ebooks; it worked!</p>
<p></p>
<br /><br />
Posted Mon, Jan 28, 2019
<h1><a href="/posts/learning-day-1-golang/">Learning Day 1: go</a></h1>
<p>This is the first Learning Day Log I&rsquo;m publishing, and it&rsquo;s about Go.</p>
<p></p>
<br /><br />
Posted Sat, Jan 26, 2019
<h1><a href="/posts/go-interfaces/">Go Interfaces</a></h1>
<p>I did some work recently that depended on Go interfaces and I found it both
straightforward and elegant.</p>
<p></p>
<br /><br />
Posted Wed, Jan 23, 2019
<h1><a href="/posts/the-evolution-of-minotaur/">The Evolution of The Minotaur</a></h1>
<p>I have a tool called The Minotaur that I just rewrote for the third time, and I
think, maybe, it&rsquo;s done.</p>
<p></p>
<br /><br />
Posted Mon, Jan 14, 2019
<h1><a href="/posts/self-control-on-a-phone/">Self-Control on a Phone</a></h1>
<p>Today I discovered that a lot of people feel alone in how they feel chained, in
one way or another, to their phones. I started the fight against that recently
and thought my findings might help other people.</p>
<p></p>
<br /><br />
Posted Thu, Jan 10, 2019
<h1><a href="/posts/notes-linking-update/">Updates to my Notes Linking Tools</a></h1>
<p>I recently improved some of my notes tools, most especially around linking to
emails.</p>
<p></p>
<br /><br />
Posted Tue, Jan 8, 2019
<h1><a href="/posts/goals-2019/">Goals for 2019</a></h1>
<p>As many do, I am attempting to affect 2019 by picking skills to improve,
subjects to learn, ways I hope to improve as a person, and then deriving
(hopefully) concrete milestones to benchmark that progress.</p>
<p></p>
<br /><br />
Posted Sun, Dec 30, 2018

<h1><a href="/posts/golang-self-signed-and-pinned-certs/">Self-Signed and Pinned Certificates in Go</a></h1>
<p>I recently needed to generate some TLS certificates in Go and trust them.
Here&rsquo;s how I did it.</p>
<p></p>
<br /><br />
Posted Sun, Dec 23, 2018
<h1><a href="/posts/validating-kubernetes-manifests/">Validating Kubernetes Manifests</a></h1>
<p>At <a href="https://www.ziprecruiter.com/hiring/technology">ZipRecruiter</a> my team is
hard at work making Kubernetes our production platform. This is an incredible
effort and I can only take the credit for very small parts of it. The issue
that I was tasked with most recently was to verify and transform Kubernetes
manifests; this post demonstrates how to do that reliably.</p>
<p></p>
<br /><br />
Posted Tue, Dec 18, 2018
<h1><a href="/posts/go-generate/">go generate: barely a framework</a></h1>
<p>I&rsquo;ve been leaning on <code>go generate</code> at work a lot lately and, when discussing it
with friends, found that they had trouble understanding it. I figured I&rsquo;d show
some examples to help.</p>
<p></p>
<br /><br />
Posted Mon, Nov 19, 2018
<h1><a href="/posts/golang-no-generics/">Go Doesn&#39;t Have Generics</a></h1>
<p>Go doesn&rsquo;t have generics. This isn&rsquo;t news, but it&rsquo;s more foundational than many
might realize.</p>
<p></p>
<br /><br />
Posted Mon, Nov 12, 2018
<h1><a href="/posts/golang-concurrency-patterns/">Go Concurrency Patterns</a></h1>
<p>I&rsquo;ve been spending some time the past couple of weeks playing with some of <a href="/posts/benefits-using-golang-adhoc-code-leatherman/">my
personal Go tools.</a> Nearly
everything I did involved concurrency, for a change. I&rsquo;ll document how I did it
and some of the wisdom I&rsquo;ve gathered from others here.</p>
<p></p>
<br /><br />
Posted Mon, Oct 22, 2018
<h1><a href="/posts/atomic-directory-population-in-golang/">Atomically Directory Population in Go</a></h1>
<p>At <a href="https://www.ziprecruiter.com/hiring/technology">work</a> I&rsquo;m building a little
tool to write data from <a href="https://docs.aws.amazon.com/secretsmanager/latest/userguide/intro.html">AWS Secrets
Manager</a>
to a directory on disk. I wrote a little package to write the secrets
atomically, because that seemed safest at the time. In retrospect just writing
each file atomically probably would have been good enough. Code and discussion
are below.</p>
<p></p>
<br /><br />
Posted Tue, Sep 18, 2018
<h1><a href="/posts/gophercon-2018/">GopherCon 2018</a></h1>
<p>This year I went to GopherCon. This post is a grab bag of what I thought was
interesting and some thoughts on this conference vs others and conferences in
general.</p>
<p></p>
<br /><br />
Posted Tue, Sep 4, 2018
<h1><a href="/posts/log-loss-detection/">Log Loss Detection</a></h1>
<p>We spent hours debugging a logging issue Friday and Monday. If you use UUIDs in
Perl, you should read this post.</p>
<p></p>
<br /><br />
Posted Wed, Jul 25, 2018
<h1><a href="/posts/some-cool-new-tools/">Some Cool New Tools</a></h1>
<p>I&rsquo;ve written (and ported) some new tools and thought others might find them
useful or inspiring.</p>
<p></p>
<br /><br />
Posted Tue, Jul 17, 2018
<h1><a href="/posts/unproductive/">unproductive</a></h1>
<p>I&rsquo;ve always wanted to carefully measure my activity on the computer and recently
<a href="https://github.com/frioux/unproductive">built a tool called <code>unproductive</code></a> to
make it happen.</p>
<p></p>
<br /><br />
Posted Thu, Jul 12, 2018
<h1><a href="/posts/announcing-shellquote/">Announcing shellquote</a></h1>
<p>In my effort to <a href="/posts/benefits-using-golang-adhoc-code-leatherman/">port certain tools to
go</a> I&rsquo;ve authored another
package: <code>github.com/frioux/shellquote</code>.</p>
<p></p>
<br /><br />
Posted Thu, Jul 5, 2018
<h1><a href="/posts/detecting-who-used-ec2-metadata-server-bcc/">Detecting who used the EC2 metadata server with BCC</a></h1>
<p>Recently at work we had a minor incident involving exhaustion of the EC2
metadata server on some of our hosts. I was able to get enough detail to
delegate the rest to a team to fix the issue.</p>
<p></p>
<br /><br />
Posted Thu, Jun 21, 2018
<h1><a href="/posts/centralized-known-hosts-for-ssh/">Centralized known_hosts for ssh</a></h1>
<p>I just wrote some code to make a (hopefully) trustworthy, shared known_hosts
file for our whole company. A handy side benefit is that it also grant us
hostname tab completion.</p>
<p></p>
<br /><br />
Posted Fri, Jun 1, 2018
<h1><a href="/posts/buffered-channels-in-golang/">Buffered Channels in Golang</a></h1>
<p>A few weeks ago when I was reading <a target="_blank" href="https://www.amazon.com/gp/product/0134190440/ref=as_li_tl?ie=UTF8&camp=1789&creative=9325&creativeASIN=0134190440&linkCode=as2&tag=afoolishmanif-20&linkId=7a70d548d8d1ab0e0baf86848938c69a">The
Go Programming Language</a><img src="//ir-na.amazon-adsystem.com/e/ir?t=afoolishmanif-20&l=am2&o=1&a=0134190440" width="1" height="1" border="0" alt="" style="border:none !important; margin:0px
!important;" /> I was reading about buffered channels and had a gut instinct
that I could write some code taking advantage of them in a precise way.
This was the comical code that came out of it.</p>
<p></p>
<br /><br />
Posted Mon, May 14, 2018
<h1><a href="/posts/c-golang-perl-and-unix/">C, Golang, Perl, and Unix</a></h1>
<p>Over the past couple months I have had the somewhat uncomfortable realization
that some of my assumptions about <em>all programs</em> are wrong. Read all about the
journey involving Unix, C, Perl, and Go.</p>
<p></p>
<br /><br />
Posted Tue, May 1, 2018
<h1><a href="/posts/announcing-mozcookiejar-golang/">Announcing mozcookiejar</a></h1>
<p>I built a little package for loading Firefox cookies into my Go tools!</p>
<p></p>
<br /><br />
Posted Fri, Apr 20, 2018
<h1><a href="/posts/reflections-on-ngs-machine-learning/">Reflections on Ng&#39;s Machine Learning</a></h1>
<p>I recently took <a href="https://www.coursera.org/learn/machine-learning/home/welcome">Andrew Ng&rsquo;s Machine Learning class on Coursera</a>; here were my
takeaways.</p>
<p></p>
<br /><br />
Posted Tue, Feb 27, 2018
<h1><a href="/posts/categorically-solving-cronspam/">Categorically Solving Cronspam</a></h1>
<p>For a little over a year at
<a href="https://www.ziprecruiter.com/hiring/technology">ZipRecruiter</a> we have had some tooling that
&ldquo;fixes&rdquo; a non-trivial amount of cronspam. Read on to see what I mean and how.</p>
<p></p>
<br /><br />
Posted Mon, Feb 26, 2018
<h1><a href="/posts/exponential-backoff-in-service-startup/">Exponential Backoff in Service Startup</a></h1>
<p>I recently added exponential backoff to service startup. Read how here.</p>
<p></p>
<br /><br />
Posted Thu, Feb 22, 2018
<h1><a href="/posts/some-code-i-deleted/">Some Code I Deleted</a></h1>
<p>I recently deleted a couple non-trivial scripts from my dotfiles and I&rsquo;m proud
of that.
</p>
<br /><br />
Posted Tue, Feb 20, 2018
<h1><a href="/posts/full-disk-whats-next/">Full Disk, What&#39;s Next?</a></h1>
<p>I recently automated yet another part of my disk usage tool. Read about it
here.</p>
<p></p>
<br /><br />
Posted Mon, Feb 19, 2018
<h1><a href="/posts/gnuplot-super-handy/">gnuplot is Super Handy</a></h1>
<p>Yesterday I wanted to graph some data by date but I didn&rsquo;t want to mess with
spreadsheet software or other graphing libraries. I reached for gnuplot after
hearing good things over the years. The results were great.</p>
<p></p>
<br /><br />
Posted Fri, Feb 16, 2018
<h1><a href="/posts/benefits-using-golang-adhoc-code-leatherman/">Benefits of using Golang for ad-hoc code: Leatherman</a></h1>
<p>I recently stumbled upon a pattern that motivates me to write little scripts in
Go instead of my normal default. I was surprised at some of the benefits.</p>
<p></p>
<br /><br />
Posted Fri, Jan 12, 2018
<h1><a href="/posts/a-love-letter-to-plain-text/">A Love Letter to Plain Text</a></h1>
<p>I have used Hugo, the blog engine this blog runs on top of, more and more lately
for less and less typical use cases. Hopefully this post will inspire others in
similar ways.</p>
<p></p>
<br /><br />
Posted Tue, Jan 2, 2018
<h1><a href="/posts/editing-registers-in-vim-regedit/">Editing Registers in Vim: RegEdit.vim</a></h1>
<p>I recently came up with the most satisfying way to edit registers in Vim I&rsquo;ve
ever seen. I hope you like it as much as I do.</p>
<p></p>
<br /><br />
Posted Fri, Oct 20, 2017
<h1><a href="/posts/advanced-projectionist-templates/">Advanced Projectionist Templates</a></h1>
<p>This week I migrated some of the vim tooling I use for <a href="/posts/hugo-unix-vim-integration/">my blog</a> from
<a href="https://github.com/sirver/ultisnips">UltiSnips</a> to <a href="https://github.com/tpope/vim-projectionist">projectionist</a>. The result is a
lighter weight and a more user friendly (for me) interface.</p>
<p></p>
<br /><br />
Posted Mon, Oct 16, 2017
<h1><a href="/posts/monitoring-service-start-stop-in-upstart/">Monitoring Service start/stop in Upstart</a></h1>
<p>Recently at ZipRecruiter I implemented a tool to ensure that we know if some
service is crashlooping. It was really easy thanks to Upstart but it took
almost a whole day to get just right.</p>
<p></p>
<br /><br />
Posted Mon, Sep 25, 2017
<h1><a href="/posts/content-based-filetype-detection-in-vim/">Content Based Filetype Detection in Vim</a></h1>
<p>Yesterday I spent a little over an hour finally figuring out how to detect a
file based on its contents in vim. It&rsquo;s pretty easy!</p>
<p></p>
<br /><br />
Posted Wed, Sep 20, 2017
<h1><a href="/posts/json-on-the-command-line/">JSON on the Command Line</a></h1>
<p>Recently my coworker Andy Ruder was complaining that he often reached for grep
when filtering JSON, and I offered to give him some tips. This post is an
expansion of what I told him.</p>
<p></p>
<br /><br />
Posted Mon, Sep 18, 2017
<h1><a href="/posts/vim-debugging/">Vim Debugging</a></h1>
<p>I use Vim quite a bit and fairly heavily, so I run into a good amount of bugs.
I&rsquo;ll share a couple tricks I&rsquo;ve learned that help debug vim.</p>
<p></p>
<br /><br />
Posted Fri, Sep 8, 2017
<h1><a href="/posts/investigation-why-sqs-slow/">Investigation: Why is SQS so slow?</a></h1>
<p>Recently I spent time figuring out why sending items to our message queue often
took absurdl100 19906    0 19906    0     0   223k      0 --:--:-- --:--:-- --:--:--  223k
y long. I am really pleased with both my solutions and my methodogy,
maybe you will be too.</p>
<p></p>
<br /><br />
Posted Sun, Aug 20, 2017
<h1><a href="/posts/supervisors-and-init-systems-7/">Supervisors and Init Systems: Part 7</a></h1>
<p>This post is the seventh in my <a href="/tags/supervisors">series about supervisors</a> and I&rsquo;m
discussing some ideas that I&rsquo;ve had while writing this series.</p>
<p></p>
<br /><br />
Posted Wed, Aug 2, 2017
<h1><a href="/posts/supervisors-and-init-systems-6/">Supervisors and Init Systems: Part 6</a></h1>
<p>This post is the sixth in my <a href="/tags/supervisors">series about supervisors</a>. I&rsquo;ll
spare you the recap since it&rsquo;s getting silly at this point. This post is about
readiness protocols.</p>
<p></p>
<br /><br />
Posted Mon, Jul 31, 2017
<h1><a href="/posts/supervisors-and-init-systems-5/">Supervisors and Init Systems: Part 5</a></h1>
<p>This post is the fifth in my <a href="/tags/supervisors">series about supervisors</a>. <a href="/posts/supervisors-and-init-systems-1/">The
first</a> <a href="/posts/supervisors-and-init-systems-2/">two posts</a> were about traditional supervisors. <a href="/posts/supervisors-and-init-systems-3/">The third</a> was
about some more unusual options. <a href="/posts/supervisors-and-init-systems-4/">The fourth</a> was about the current most
popular choices. This post is about some of the unusual trends I&rsquo;ve noticed
during my three year long obsession with supervisors.</p>
<p></p>
<br /><br />
 Posted Wed, Jul 26, 2017
</div>
</section>
</div>
</div>
<ul class="pagination">
<li class="page-item">
<a href="/" class="page-link" aria-label="First"><span aria-hidden="true">&laquo;&laquo;</span></a>
</li>
<li class="page-item disabled">
<a href="" class="page-link" aria-label="Previous"><span aria-hidden="true">&laquo;</span></a>
</li>
<li class="page-item active"><a class="page-link" href="/">1</a></li>
<li class="page-item"><a class="page-link" href="/page/2/">2</a></li>
<li class="page-item"><a class="page-link" href="/page/3/">3</a></li>
<li class="page-item disabled"><span aria-hidden="true">&nbsp;&hellip;&nbsp;</span></li>
<li class="page-item"><a class="page-link" href="/page/10/">10</a></li>
<li class="page-item">
<a href="/page/2/" class="page-link" aria-label="Next"><span aria-hidden="true">&raquo;</span></a>
</li>
<li class="page-item">
<a href="/page/10/" class="page-link" aria-label="Last"><span aria-hidden="true">&raquo;&raquo;</span></a>
</li>
</ul>
</div>
</div>
</div>
</div>
<div class="container">
<hr>
<footer id="footer">
<p class="pull-right"><a href="#top">Back to top</a></p>
<ul id="tags">
<li><a href="/tags">all tags</a></li>
</ul>
</footer>
</div>
<script src="/static/js/jquery.js" type="568c67a7cf35485aee2a5b8b-text/javascript"></script>
<script src="/static/js/bootstrap.min.js" type="568c67a7cf35485aee2a5b8b-text/javascript"></script>
<script type="568c67a7cf35485aee2a5b8b-text/javascript">
        $(document).ready(function() {
            $("nav#TableOfContents a").click(function() {
                $("html, body").animate({
                    scrollTop: $($(this).attr("href")).offset().top-25 + "px"
                }, {
                    duration: 450,
                });
                return false;
            });
        });
    </script>
<script src="https://ajax.cloudflare.com/cdn-cgi/scripts/a2bd7673/cloudflare-static/rocket-loader.min.js" data-cf-settings="568c67a7cf35485aee2a5b8b-|49" defer=""></script></body>
</html>
`)
	}))
	defer ts.Close()

	buf := &bytes.Buffer{}
	err := run(strings.NewReader(ts.URL), buf)
	assert.NoError(t, err)
	assert.Equal(t, "[fREW Schmidt's Foolish Manifesto]("+ts.URL+")\n", buf.String())
}
