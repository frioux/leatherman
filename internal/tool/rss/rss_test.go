package rss

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Parallel()

	f, err := ioutil.TempFile("", "*.js")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintln(w, `
<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>fREW Schmidt&#39;s Foolish Manifesto on fREW Schmidt&#39;s Foolish Manifesto</title>
    <link>https://blog.afoolishmanifesto.com/</link>
    <description>Recent content in fREW Schmidt&#39;s Foolish Manifesto on fREW Schmidt&#39;s Foolish Manifesto</description>
    <generator>Hugo -- gohugo.io</generator>
    <lastBuildDate>Thu, 21 Mar 2019 07:25:18 +0000</lastBuildDate>
    <atom:link href="/" rel="self" type="application/rss+xml" />

    <item>
      <title>Sorting Books</title>
      <link>https://blog.afoolishmanifesto.com/posts/sorting-books/</link>
      <pubDate>Thu, 21 Mar 2019 07:25:18 +0000</pubDate>

      <guid isPermaLink="false">18e35dc0-5e01-4dd2-af7a-9a273134203f</guid>
      <description>&lt;p&gt;I wrote a little program to sort lists of books.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>perl</category>

      <category>frew-warez</category>

    </item>

    <item>
      <title>Automating Email</title>
      <link>https://blog.afoolishmanifesto.com/posts/automating-email/</link>
      <pubDate>Mon, 18 Mar 2019 07:10:42 +0000</pubDate>

      <guid isPermaLink="false">ddbf4a02-d7b1-4736-8f0d-b5693027a6ca</guid>
      <description>&lt;p&gt;I just automated a couple common email tasks.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>mutt</category>

      <category>golang</category>

    </item>

    <item>
      <title>How to Add a Subscription Service to Your Blog</title>
      <link>https://blog.afoolishmanifesto.com/posts/how-to-add-a-subscription-mode-to-your-blog/</link>
      <pubDate>Thu, 07 Mar 2019 07:15:57 +0000</pubDate>

      <guid isPermaLink="false">0cf2f92a-232c-4b25-a2f7-48dedb0e723b</guid>
      <description>&lt;p&gt;I used to use a service to email subscribers updates to my blog.  The service
broke, but I automated my way around it.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>perl</category>

      <category>blog</category>

      <category>meta</category>

    </item>

    <item>
      <title>Fixing Buggy Haskell Programs with Go</title>
      <link>https://blog.afoolishmanifesto.com/posts/fixing-buggy-haskell-programs-with-golang/</link>
      <pubDate>Wed, 27 Feb 2019 07:11:08 +0000</pubDate>

      <guid isPermaLink="false">b940dc2a-6ebd-4a0f-b6c2-3a5f452e2230</guid>
      <description>&lt;p&gt;I recently ran into a stupid bug in a program written in Haskell and found it
much easier to paper over with a few lines of Go than to properly fix.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>haskell</category>

      <category>golang</category>

    </item>

    <item>
      <title>Learning Day 2: DIY Games</title>
      <link>https://blog.afoolishmanifesto.com/posts/learning-day-2-diy-games/</link>
      <pubDate>Sat, 23 Feb 2019 19:41:55 +0000</pubDate>

      <guid isPermaLink="false">360a79f0-5b2f-48e0-a5e1-3e0e79d000e0</guid>
      <description>&lt;p&gt;Today I did my second Learning Day; the subject was DIY Games.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>lua</category>

      <category>pico-8</category>

      <category>learning-day</category>

      <category>self</category>

    </item>

    <item>
      <title>Busting the Cloudflare Cache</title>
      <link>https://blog.afoolishmanifesto.com/posts/busting-cloudflare-cache/</link>
      <pubDate>Wed, 20 Feb 2019 07:15:17 +0000</pubDate>

      <guid isPermaLink="false">96139cd2-b350-4d4e-9a6e-045645ba8cdd</guid>
      <description>&lt;p&gt;I automated blowing the cache for this blog.  Read on to see how I did it.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>perl</category>

      <category>meta</category>

      <category>cloudflare</category>

    </item>

    <item>
      <title>graphviz describing multi-stage docker builds</title>
      <link>https://blog.afoolishmanifesto.com/posts/graphviz/</link>
      <pubDate>Mon, 11 Feb 2019 07:27:10 +0000</pubDate>

      <guid isPermaLink="false">f35be163-f9b1-475b-b4c5-abc0d149bc6f</guid>
      <description>&lt;p&gt;I recently decided I should learn to use Graphviz more, as a great tool for
making certain kinds of plots.  Less than a week later a great use case
surfaced.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>tool</category>

      <category>graphviz</category>

      <category>docker</category>

      <category>ziprecruiter</category>

      <category>perl</category>

    </item>

    <item>
      <title>Amygdala</title>
      <link>https://blog.afoolishmanifesto.com/posts/amygdala/</link>
      <pubDate>Tue, 05 Feb 2019 07:12:26 +0000</pubDate>

      <guid isPermaLink="false">bca651f1-8ba4-4f18-9efe-b4b869f7bedc</guid>
      <description>&lt;p&gt;This past weekend I started re-creating a tool I used to have, using new tools,
techniques, and infrastructure.  The tool allows, at least, adding to my own
todo list via SMS.  It&amp;rsquo;s working great!&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>golang</category>

      <category>perl</category>

      <category>amygdala</category>

    </item>

    <item>
      <title>Deploying to Kubernetes at ZipRecruiter</title>
      <link>https://blog.afoolishmanifesto.com/posts/deploying-to-kubernetes-at-ziprecruiter/</link>
      <pubDate>Wed, 30 Jan 2019 07:36:37 +0000</pubDate>

      <guid isPermaLink="false">fcc31a7f-2696-45a8-8585-bbbf9ce521d6</guid>
      <description>&lt;p&gt;At &lt;a href=&#34;https://www.ziprecruiter.com/hiring/technology&#34;&gt;ZR&lt;/a&gt; we are working hard to
get stuff migrated to Kubernetes, and a big part of that is our cicd pipeline.
We have that stable enough that I can explain the major parts.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>ziprecruiter</category>

      <category>kubernetes</category>

      <category>cicd</category>

    </item>

    <item>
      <title>Full Text Search for ebooks</title>
      <link>https://blog.afoolishmanifesto.com/posts/full-text-search-for-ebooks/</link>
      <pubDate>Mon, 28 Jan 2019 07:30:26 +0000</pubDate>

      <guid isPermaLink="false">78bcacf7-dc50-4fdf-91c5-8365ab61c86f</guid>
      <description>&lt;p&gt;This past weekend I did a learning day that inspired me to try SQLite for
indexing my ebooks; it worked!&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>meta</category>

      <category>learning-day</category>

    </item>

    <item>
      <title>Learning Day 1: go</title>
      <link>https://blog.afoolishmanifesto.com/posts/learning-day-1-golang/</link>
      <pubDate>Sat, 26 Jan 2019 16:46:28 +0000</pubDate>

      <guid isPermaLink="false">2122f364-8a42-4734-880e-c5da312b7a5e</guid>
      <description>&lt;p&gt;This is the first Learning Day Log I&amp;rsquo;m publishing, and it&amp;rsquo;s about Go.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>golang</category>

      <category>learning-day</category>

      <category>meta</category>

    </item>

    <item>
      <title>Go Interfaces</title>
      <link>https://blog.afoolishmanifesto.com/posts/go-interfaces/</link>
      <pubDate>Wed, 23 Jan 2019 08:30:03 +0000</pubDate>

      <guid isPermaLink="false">7a23bd20-d454-4384-bf0e-b5ccddf85833</guid>
      <description>&lt;p&gt;I did some work recently that depended on Go interfaces and I found it both
straightforward and elegant.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>golang</category>

      <category>programming</category>

      <category>programming-languages</category>

    </item>

    <item>
      <title>The Evolution of The Minotaur</title>
      <link>https://blog.afoolishmanifesto.com/posts/the-evolution-of-minotaur/</link>
      <pubDate>Mon, 14 Jan 2019 07:33:50 +0000</pubDate>

      <guid isPermaLink="false">4e448322-1f08-4749-b8c2-607aac3dd5e4</guid>
      <description>&lt;p&gt;I have a tool called The Minotaur that I just rewrote for the third time, and I
think, maybe, it&amp;rsquo;s done.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>perl</category>

      <category>golang</category>

      <category>ziprecruiter</category>

      <category>mitsi</category>

      <category>meta</category>

      <category>toolsmith</category>

    </item>

    <item>
      <title>Self-Control on a Phone</title>
      <link>https://blog.afoolishmanifesto.com/posts/self-control-on-a-phone/</link>
      <pubDate>Thu, 10 Jan 2019 19:28:00 +0000</pubDate>

      <guid isPermaLink="false">0d510493-61a5-4fb8-b93f-29f570befd77</guid>
      <description>&lt;p&gt;Today I discovered that a lot of people feel alone in how they feel chained, in
one way or another, to their phones.  I started the fight against that recently
and thought my findings might help other people.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>meta</category>

      <category>self-control</category>

      <category>super-powers</category>

      <category>phone</category>

    </item>

    <item>
      <title>Updates to my Notes Linking Tools</title>
      <link>https://blog.afoolishmanifesto.com/posts/notes-linking-update/</link>
      <pubDate>Tue, 08 Jan 2019 08:11:00 +0000</pubDate>

      <guid isPermaLink="false">2d7780f0-d095-4df6-97ac-cc1802b44cf5</guid>
      <description>&lt;p&gt;I recently improved some of my notes tools, most especially around linking to
emails.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>frew-warez</category>

      <category>golang</category>

      <category>meta</category>

      <category>vim</category>

    </item>

    <item>
      <title>Goals for 2019</title>
      <link>https://blog.afoolishmanifesto.com/posts/goals-2019/</link>
      <pubDate>Sun, 30 Dec 2018 08:10:28 +0000</pubDate>

      <guid isPermaLink="false">bd3d010e-b286-4903-8d54-f8844a591cb4</guid>
      <description>&lt;p&gt;As many do, I am attempting to affect 2019 by picking skills to improve,
subjects to learn, ways I hope to improve as a person, and then deriving
(hopefully) concrete milestones to benchmark that progress.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>goals</category>

      <category>meta</category>

    </item>

    <item>
      <title>Self-Signed and Pinned Certificates in Go</title>
      <link>https://blog.afoolishmanifesto.com/posts/golang-self-signed-and-pinned-certs/</link>
      <pubDate>Sun, 23 Dec 2018 07:29:05 +0000</pubDate>

      <guid isPermaLink="false">4e8b5670-3908-4ced-9ce7-b0f5dabfe085</guid>
      <description>&lt;p&gt;I recently needed to generate some TLS certificates in Go and trust them.
Here&amp;rsquo;s how I did it.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>golang</category>

      <category>ssl</category>

      <category>tls</category>

    </item>

    <item>
      <title>Validating Kubernetes Manifests</title>
      <link>https://blog.afoolishmanifesto.com/posts/validating-kubernetes-manifests/</link>
      <pubDate>Tue, 18 Dec 2018 07:20:15 +0000</pubDate>

      <guid isPermaLink="false">0d291e43-0f72-4922-8790-275a114c951e</guid>
      <description>&lt;p&gt;At &lt;a href=&#34;https://www.ziprecruiter.com/hiring/technology&#34;&gt;ZipRecruiter&lt;/a&gt; my team is
hard at work making Kubernetes our production platform.  This is an incredible
effort and I can only take the credit for very small parts of it.  The issue
that I was tasked with most recently was to verify and transform Kubernetes
manifests; this post demonstrates how to do that reliably.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>kubernetes</category>

      <category>perl</category>

      <category>golang</category>

    </item>

    <item>
      <title>go generate: barely a framework</title>
      <link>https://blog.afoolishmanifesto.com/posts/go-generate/</link>
      <pubDate>Mon, 19 Nov 2018 07:20:59 +0000</pubDate>

      <guid isPermaLink="false">fd338831-8f40-4b03-8bf6-144833a1112d</guid>
      <description>&lt;p&gt;I&amp;rsquo;ve been leaning on &lt;code&gt;go generate&lt;/code&gt; at work a lot lately and, when discussing it
with friends, found that they had trouble understanding it.  I figured I&amp;rsquo;d show
some examples to help.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>golang</category>

    </item>

    <item>
      <title>Go Doesn&#39;t Have Generics</title>
      <link>https://blog.afoolishmanifesto.com/posts/golang-no-generics/</link>
      <pubDate>Mon, 12 Nov 2018 09:37:49 +0000</pubDate>

      <guid isPermaLink="false">602effcf-b9e9-4e13-afb8-4a08907b3ead</guid>
      <description>&lt;p&gt;Go doesn&amp;rsquo;t have generics.  This isn&amp;rsquo;t news, but it&amp;rsquo;s more foundational than many
might realize.&lt;/p&gt;

&lt;p&gt;&lt;/p&gt;</description>

      <category>golang</category>

      <category>psa</category>

    </item>

  </channel>
</rss>
		`)
	}))
	defer ts.Close()

	_, err = f.WriteString("[]")
	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	err = run(ts.URL, f.Name(), buf)
	assert.NoError(t, err)
	assert.Equal(t, buf.String(), `[Sorting Books](https://blog.afoolishmanifesto.com/posts/sorting-books/)
[Automating Email](https://blog.afoolishmanifesto.com/posts/automating-email/)
[How to Add a Subscription Service to Your Blog](https://blog.afoolishmanifesto.com/posts/how-to-add-a-subscription-mode-to-your-blog/)
[Fixing Buggy Haskell Programs with Go](https://blog.afoolishmanifesto.com/posts/fixing-buggy-haskell-programs-with-golang/)
[Learning Day 2: DIY Games](https://blog.afoolishmanifesto.com/posts/learning-day-2-diy-games/)
[Busting the Cloudflare Cache](https://blog.afoolishmanifesto.com/posts/busting-cloudflare-cache/)
[graphviz describing multi-stage docker builds](https://blog.afoolishmanifesto.com/posts/graphviz/)
[Amygdala](https://blog.afoolishmanifesto.com/posts/amygdala/)
[Deploying to Kubernetes at ZipRecruiter](https://blog.afoolishmanifesto.com/posts/deploying-to-kubernetes-at-ziprecruiter/)
[Full Text Search for ebooks](https://blog.afoolishmanifesto.com/posts/full-text-search-for-ebooks/)
[Learning Day 1: go](https://blog.afoolishmanifesto.com/posts/learning-day-1-golang/)
[Go Interfaces](https://blog.afoolishmanifesto.com/posts/go-interfaces/)
[The Evolution of The Minotaur](https://blog.afoolishmanifesto.com/posts/the-evolution-of-minotaur/)
[Self-Control on a Phone](https://blog.afoolishmanifesto.com/posts/self-control-on-a-phone/)
[Updates to my Notes Linking Tools](https://blog.afoolishmanifesto.com/posts/notes-linking-update/)
[Goals for 2019](https://blog.afoolishmanifesto.com/posts/goals-2019/)
[Self-Signed and Pinned Certificates in Go](https://blog.afoolishmanifesto.com/posts/golang-self-signed-and-pinned-certs/)
[Validating Kubernetes Manifests](https://blog.afoolishmanifesto.com/posts/validating-kubernetes-manifests/)
[go generate: barely a framework](https://blog.afoolishmanifesto.com/posts/go-generate/)
[Go Doesn't Have Generics](https://blog.afoolishmanifesto.com/posts/golang-no-generics/)
`)
}
