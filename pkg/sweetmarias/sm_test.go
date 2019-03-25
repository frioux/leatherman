package sweetmarias

import (
	"fmt"
	"net/http"
	"testing"

	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func TestLoadCoffee(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `
<!doctype html>
<html lang="en-US">
    <head prefix="og: http://ogp.me/ns# fb: http://ogp.me/ns/fb# product: http://ogp.me/ns/product#">
        <script>
    var require = {
        "baseUrl": "https://sweet-marias-live-sweetmarias.netdna-ssl.com/static/version1553065162/frontend/Forix/sweetmarias/en_US"
    };
</script>
        <meta charset="utf-8"/>
<meta name="description" content="Honey process seems to bring out fruited notes like cranberry, raisin, red grape, and underscored by molasses-like sweetness. This PNG boasts body, and with mild acidity, is great espresso too. City+ to Full City+. Good for espresso."/>
<meta name="keywords" content="Green Coffee"/>
<meta name="robots" content="INDEX,FOLLOW"/>
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1.0, user-scalable=no"/>
<title>Papua New Guinea Honey Process Nebilyer Estate</title>
<link  rel="stylesheet" type="text/css"  media="all" href="https://sweet-marias-live-sweetmarias.netdna-ssl.com/static/version1553065162/_cache/merged/7d2d69c418360df7c8d3c84701d181c0.min.css" />
<link  rel="stylesheet" type="text/css"  media="screen and (min-width: 768px)" href="https://sweet-marias-live-sweetmarias.netdna-ssl.com/static/version1553065162/frontend/Forix/sweetmarias/en_US/css/styles-l.min.css" />
<link  rel="stylesheet" type="text/css"  media="print" href="https://sweet-marias-live-sweetmarias.netdna-ssl.com/static/version1553065162/frontend/Forix/sweetmarias/en_US/css/print.min.css" />
<script  type="text/javascript"  src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/static/version1553065162/_cache/merged/54cf402f9c56edcceaedb479dec48d23.min.js"></script>
<link  rel="stylesheet" type="text/css" href="https://fonts.googleapis.com/css?family=Carrois+Gothic+SC|Roboto+Condensed:400,700|Roboto+Mono:300,400,500,700" />
<link  rel="icon" type="image/x-icon" href="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/favicon/stores/1/Favicon48_sweet.png" />
<link  rel="shortcut icon" type="image/x-icon" href="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/favicon/stores/1/Favicon48_sweet.png" />
<link  rel="canonical" href="https://www.sweetmarias.com/papua-new-guinea-honey-process-nebilyer-estate.html" />
<!--Start of Zendesk Chat Script-->
<script type="text/javascript">
window.$zopim||(function(d,s){var z=$zopim=function(c){z._.push(c)},$=z.s=
d.createElement(s),e=d.getElementsByTagName(s)[0];z.set=function(o){z.set.
_.push(o)};z._=[];z.set._=[];$.async=!0;$.setAttribute("charset","utf-8");
$.src="https://v2.zopim.com/?5TLWa8qii2SFYKHBC2a6eQGHWu6U28a2";z.t=+new Date;$.
type="text/javascript";e.parentNode.insertBefore($,e)})(document,"script");
</script>
<!--End of Zendesk Chat Script-->        <script type="text/javascript" src="https://chimpstatic.com/mcjs-connected/js/users/ece32b8792ad4351762abfe02/cde1fc1c2e704cc4e380fbe08.js" async></script>
<link rel="stylesheet" type="text/css" media="all"
      href="//maxcdn.bootstrapcdn.com/font-awesome/latest/css/font-awesome.min.css"/>
<meta property="og:type" content="og:product" />
<meta property="og:title" content="Papua New Guinea Honey Nebilyer Estate" />
<meta property="og:image" content="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/catalog/product/cache/image/265x265/beff4985b56e3afdbeabfc89641a4582/p/a/papua-new-guinea-sweetmarias-2.jpg" />
<meta property="og:description" content="Honey process seems to bring out fruited notes like cranberry, raisin, red grape, and underscored by molasses-like sweetness. This PNG boasts body, and with mild acidity, is great espresso too. City+ to Full City+. Good for espresso." />
<meta property="og:url" content="https://www.sweetmarias.com/green-coffee/oceania/papua-new-guinea/papua-new-guinea-honey-process-nebilyer-estate.html" />
    <meta property="product:price:amount" content="12.07"/>
    <meta property="product:price:currency" content="USD"/>
    </head>
    <body data-container="body" data-mage-init='{"loaderAjax": {}, "loader": { "icon": "https://sweet-marias-live-sweetmarias.netdna-ssl.com/static/version1553065162/frontend/Forix/sweetmarias/en_US/images/loader-2.gif"}}' itemtype="http://schema.org/Product" itemscope="itemscope" class="page-product-configurable catalog-product-view product-papua-new-guinea-honey-process-nebilyer-estate categorypath-green-coffee-oceania-papua-new-guinea category-papua-new-guinea page-layout-1column">
            <script>
        require.config({
            deps: [
                'jquery',
                'mage/translate',
                'jquery/jquery-storageapi'
            ],
            callback: function ($) {
                'use strict';

                var dependencies = [],
                    versionObj;

                $.initNamespaceStorage('mage-translation-storage');
                $.initNamespaceStorage('mage-translation-file-version');
                versionObj = $.localStorage.get('mage-translation-file-version');

                if (versionObj.version !== '9d8a871a7c095f82e4a025917a37e7393fdde3ff') {
                    dependencies.push(
                        'text!js-translation.json'
                    );

                }

                require.config({
                    deps: dependencies,
                    callback: function (string) {
                        if (typeof string === 'string') {
                            $.mage.translate.add(JSON.parse(string));
                            $.localStorage.set('mage-translation-storage', string);
                            $.localStorage.set(
                                'mage-translation-file-version',
                                {
                                    version: '9d8a871a7c095f82e4a025917a37e7393fdde3ff'
                                }
                            );
                        } else {
                            $.mage.translate.add($.localStorage.get('mage-translation-storage'));
                        }
                    }
                });
            }
        });
    </script>

<script type="text/x-magento-init">
    {
        "*": {
            "mage/cookies": {
                "expires": null,
                "path": "/",
                "domain": ".www.sweetmarias.com",
                "secure": false,
                "lifetime": "180000"
            }
        }
    }
</script>
    <noscript>
        <div class="message global noscript">
            <div class="content">
                <p>
                    <strong>JavaScript seems to be disabled in your browser.</strong>
                    <span>For the best experience on our site, be sure to turn on Javascript in your browser.</span>
                </p>
            </div>
        </div>
    </noscript>
                <!-- BEGIN GOOGLE UNIVERSAL ANALYTICS CODE -->
                <script>
                    //<![CDATA[
                    (function (i, s, o, g, r, a, m) {
                        i['GoogleAnalyticsObject'] = r;
                        i[r] = i[r] || function () {
                            (i[r].q = i[r].q || []).push(arguments)
                        }, i[r].l = 1 * new Date();
                        a = s.createElement(o),
                            m = s.getElementsByTagName(o)[0];
                        a.async = 1;
                        a.src = g;
                        m.parentNode.insertBefore(a, m)
                    })(window, document, 'script', '//www.google-analytics.com/analytics.js', 'ga');

                    
ga('create', 'UA-11688411-2', 'auto');
ga('send', 'pageview');
                                        //]]>
                </script>
                <!-- END GOOGLE UNIVERSAL ANALYTICS CODE -->
            <script>
    var dlCurrencyCode = dlCurrencyCode || '';
    var dataLayer = dataLayer || [];
    var staticImpressions = staticImpressions || [];
    var staticPromotions = staticPromotions || [];
    var updatedImpressions = updatedImpressions || [];
    var updatedPromotions = updatedPromotions || [];
    var cookieAddToCart = 'add_to_cart';
    var cookieRemoveFromCart = cookieRemoveFromCart || 'remove_from_cart';
    var bannerCounter = bannerCounter || 0;

    require([
        "jquery",
        "prototype",
        "Magento_Customer/js/customer-data"
    ], function(jQuery, prototype, customerData){

        function GoogleAnalyticsUniversal(){}
        GoogleAnalyticsUniversal.prototype = {
            activeOnCategory : function(id, name, category, list, position) {
                dataLayer.push({
                    'event': 'productClick',
                    'ecommerce': {
                        'click': {
                            'actionField': {
                                'list': list
                            },
                            'products': [{
                                'id': id,
                                'name': name,
                                'category': category,
                                'list': list,
                                'position': position
                            }]
                        }
                    }
                });
            },
            activeOnProducts : function(id, name, list, position) {
                dataLayer.push({
                    'event': 'productClick',
                    'ecommerce': {
                        'click': {
                            'actionField': {
                                'list': list
                            },
                            'products': [{
                                'id': id,
                                'name': name,
                                'list': list,
                                'position': position
                            }]
                        }
                    }
                });
            },
            addToCart : function(id, name, price, quantity) {
                dataLayer.push({
                    'event': 'addToCart',
                    'ecommerce': {
                        'currencyCode' : dlCurrencyCode,
                        'add': {
                            'products': [{
                                'id': id,
                                'name': name,
                                'price': price,
                                'quantity': quantity
                            }]
                        }
                    }
                });
            },
            clickBanner : function(id, name, creative, position) {
                dataLayer.push({
                    'event': 'promotionClick',
                    'ecommerce': {
                        'promoClick': {
                            'promotions': [{
                                'id': id,
                                'name': name,
                                'creative': creative,
                                'position': position
                            }]
                        }
                    }
                });
            },
            bindImpressionClick : function(id, type, name, category, list, position, blockType, listPosition) {
                var productLink = [];
                var eventBlock;
                switch (blockType)  {
                    case 'catalog.product.related':
                        eventBlock = '.products-related .products';
                        break;
                    case 'product.info.upsell':
                        eventBlock = '.products-upsell .products';
                        break;
                    case 'checkout.cart.crosssell':
                        eventBlock = '.products-crosssell .products';
                        break;
                    case 'category.products.list':
                    case 'search_result_list':
                        eventBlock = '.products .products';
                        break;
                }
                productLink = $$(eventBlock + ' .item:nth(' + (listPosition) + ') a');
                if (type == 'configurable' || type == 'bundle' || type == 'grouped') {
                    productLink = $$(
                        eventBlock + ' .item:nth(' + (listPosition) + ') .tocart',
                        eventBlock + ' .item:nth(' + (listPosition) + ') a'
                    );
                }
                productLink.each(function(element) {
                    element.observe('click', function(event) {
                        googleAnalyticsUniversal.activeOnProducts(
                            id,
                            name,
                            list,
                            position);
                    });
                });
            },

            updateImpressions: function() {
                var pageImpressions = this.mergeImpressions();
                var dlImpressions = {
                    'event' : 'productImpression',
                    'ecommerce' : {
                        'impressions' : []
                    }
                };
                var impressionCounter = 0;
                for (blockName in pageImpressions) {
                    if (blockName === 'length' || !pageImpressions.hasOwnProperty(blockName))
                        continue;
                    for (var i = 0; i < pageImpressions[blockName].length; i++) {
                        var impression = pageImpressions[blockName][i];
                        dlImpressions.ecommerce.impressions.push({
                            'id': impression.id,
                            'name': impression.name,
                            'category': impression.category,
                            'list': impression.list,
                            'position': impression.position
                        });
                        impressionCounter++;
                        this.bindImpressionClick(
                            impression.id,
                            impression.type,
                            impression.name,
                            impression.category,
                            impression.list,
                            impression.position,
                            blockName,
                            impression.listPosition
                        );
                    }
                }
                if (impressionCounter > 0) {
                    dataLayer.push(dlImpressions);
                }
            },

            mergeImpressions: function() {
                var pageImpressions = [];
                var blockNames = ["category.products.list","product.info.upsell","catalog.product.related","checkout.cart.crosssell","search_result_list"];
                blockNames.each(function(blockName) {
                    // check if there is a new block generated by FPC placeholder update
                    if (blockName in updatedImpressions) {
                        pageImpressions[blockName] = updatedImpressions[blockName];
                    } else if (blockName in staticImpressions) { // use the static data otherwise
                        pageImpressions[blockName] = staticImpressions[blockName];
                    }
                });
                return pageImpressions;
            },

            updatePromotions : function() {
                var dlPromotions = {
                    'event' : 'promotionView',
                    'ecommerce': {
                        'promoView': {
                            'promotions' : []
                        }
                    }
                };
                var pagePromotions = [];
                // check if there is a new block generated by FPC placeholder update
                if (updatedPromotions.length) {
                    pagePromotions = updatedPromotions;
                }
                // use the static data otherwise
                if (pagePromotions.length == 0 && staticPromotions.length) {
                    pagePromotions = staticPromotions;
                }
                var promotionCounter = 0;
                var bannerIds = [];
                if (jQuery('[data-banner-id]').length) {
                    _.each(jQuery('[data-banner-id]'), function (banner) {
                        var banner = jQuery(banner);
                        var ids = (banner.data('ids') + '').split(',');
                        bannerIds = jQuery.merge(bannerIds, ids);
                    });
                }
                bannerIds = jQuery.unique(bannerIds);
                for (var i = 0; i < pagePromotions.length; i++) {
                    var promotion = pagePromotions[i];
                    if (jQuery.inArray(promotion.id, bannerIds) == -1 || promotion.activated == '0') {
                        continue;
                    }
                    dlPromotions.ecommerce.promoView.promotions.push({
                        'id': promotion.id,
                        'name': promotion.name,
                        'creative': promotion.creative,
                        'position': promotion.position
                    });
                    promotionCounter++;
                }
                if (promotionCounter > 0) {
                    dataLayer.push(dlPromotions);
                }
                jQuery('[data-banner-id]').on('click', '[data-banner-id]', function(e){
                    var bannerId = jQuery(this).attr('data-banner-id');
                    var promotions = _.filter(pagePromotions, function(item) {
                        return item.id === bannerId;
                    });
                    _.each(promotions, function(promotion) {
                        googleAnalyticsUniversal.clickBanner(
                            promotion.id,
                            promotion.name,
                            promotion.creative,
                            promotion.position
                        );
                    });
                });
            }
        };

        GoogleAnalyticsUniversalCart = function(){
            this.productQtys = [];
            this.origProducts = {};
            this.productWithChanges = [];
            this.addedProducts = [];
            this.removedProducts = [];
            this.googleAnalyticsUniversalData = {};
        };
        GoogleAnalyticsUniversalCart.prototype = {
            // ------------------- shopping cart ------------------------
            listenMinicartReload : function() {
                var context = this;
                if (typeof(Minicart) != 'undefined' && typeof(Minicart.prototype.initAfterEvents)) {
                    Minicart.prototype.initAfterEvents['GoogleAnalyticsUniversalCart:subscribeProductsUpdateInCart']
                        = function() {
                        context.subscribeProductsUpdateInCart();
                        context.parseAddToCartCookies();
                        context.parseRemoveFromCartCookies();
                    };
                    // if we are removing last item init don't calling
                    Minicart.prototype.removeItemAfterEvents[
                        'GoogleAnalyticsUniversalCart:subscribeProductsRemoveFromCart'
                        ] = function() {
                        context.parseRemoveFromCartCookies();
                    };
                }
            },
            subscribeProductsUpdateInCart : function() {
                var context = this;
                $$('[data-cart-item-update]').each(function(element) {
                    $(element).stopObserving('click').observe('click', function(){
                        context.updateCartObserver();
                    });
                });
                jQuery('[data-block="minicart"]').on('mousedown', '.update-cart-item', function(){
                    context.updateCartObserver();
                });

                $$('[data-multiship-item-update]').each(function(element) {
                    $(element).stopObserving('click').observe('click', function(){
                        context.updateMulticartCartObserver();
                    });
                });
                $$('[data-cart-empty]').each(function(element){
                    $(element).stopObserving('click').observe('click', function(){
                        context.emptyCartObserver();
                    });
                });
            },
            emptyCartObserver : function() {
                this.collectOriginalProducts();
                for (var i in this.origProducts) {
                    if (i != 'length' && this.origProducts.hasOwnProperty(i)) {
                        var product = Object.extend({}, this.origProducts[i]);
                        this.removedProducts.push(product);
                    }
                }
                this.cartItemRemoved();
            },
            updateMulticartCartObserver : function() {
                this.collectMultiProductsWithChanges();
                this.collectProductsForMessages();
                this.cartItemAdded();
                this.cartItemRemoved();
            },
            updateCartObserver : function() {
                this.collectProductsWithChanges();
                this.collectProductsForMessages();
                this.cartItemAdded();
                this.cartItemRemoved();
            },
            collectMultiProductsWithChanges : function() {
                this.collectOriginalProducts();
                this.collectMultiCartQtys();
                this.productWithChanges = [];
                var groupedProducts = {};
                for (var i = 0; i < this.productQtys.length; i++) {
                    var cartProduct = this.productQtys[i];
                    if (typeof(groupedProducts[cartProduct.id]) == 'undefined') {
                        groupedProducts[cartProduct.id] = parseInt(cartProduct.qty, 10);
                    } else {
                        groupedProducts[cartProduct.id] += parseInt(cartProduct.qty, 10);
                    }
                }
                for (var j in groupedProducts) {
                    if (groupedProducts.hasOwnProperty(j)) {
                        if (
                            (typeof(this.origProducts[j]) != 'undefined')
                            && (groupedProducts[j] != this.origProducts[j].qty)
                        ) {
                            var product = Object.extend({}, this.origProducts[j]);
                            product['qty'] = groupedProducts[j];
                            this.productWithChanges.push(product);
                        }
                    }
                }
            },
            collectProductsWithChanges : function () {
                this.collectOriginalProducts();
                this.collectCartQtys();
                this.collectMiniCartQtys();
                this.productWithChanges = [];
                for (var i = 0; i < this.productQtys.length; i++) {
                    var cartProduct = this.productQtys[i];
                    if (
                        (typeof(this.origProducts[cartProduct.id]) != 'undefined') &&
                        (cartProduct.qty != this.origProducts[cartProduct.id].qty)
                    ) {
                        var product = Object.extend({}, this.origProducts[cartProduct.id]);
                        if (parseInt(cartProduct.qty, 10) > 0) {
                            product['qty'] = cartProduct.qty;
                            this.productWithChanges.push(product);
                        }
                    }
                }
            },
            collectOriginalProducts : function() {
                var products = {};
                var items = customerData.get('cart')().items;
                if (items !== undefined) {
                    items.each(function(item) {
                        products[item.product_sku] = {
                            "id": item.product_sku,
                            "name": item.product_name,
                            "price": item.product_price_value,
                            "qty": parseInt(item.qty, 10)
                        };
                    });
                }
                this.googleAnalyticsUniversalData['shoppingCartContent'] = products;
                this.origProducts = this.googleAnalyticsUniversalData['shoppingCartContent'];
            },
            collectMultiCartQtys : function() {
                var productQtys = [];
                $$('[data-multiship-item-id]').each(function(element){
                    productQtys.push({
                        'id' : $(element).readAttribute('data-multiship-item-id'),
                        'qty' : $(element).getValue()
                    });
                });
                this.productQtys = productQtys;
            },
            collectCartQtys : function() {
                var productQtys = [];
                $$('[data-cart-item-id]').each(function(element){
                    productQtys.push({
                        'id' : $(element).readAttribute('data-cart-item-id'),
                        'qty' : $(element).getValue()
                    });
                });
                this.productQtys = productQtys;
            },
            collectMiniCartQtys : function() {
                var productQtys = [];
                $$('input[data-cart-item-id]').each(function(element){
                    productQtys.push({
                        'id' : $(element).readAttribute('data-cart-item-id'),
                        'qty' : $(element).getValue()
                    });
                });
                this.productQtys = productQtys;
            },
            collectProductsForMessages : function() {
                this.addedProducts = [];
                this.removedProducts = [];
                for (var i = 0; i < this.productWithChanges.length; i++) {
                    var product = this.productWithChanges[i];
                    if (typeof(this.origProducts[product.id]) != 'undefined') {
                        if (product.qty > this.origProducts[product.id].qty) {
                            product.qty = Math.abs(product.qty - this.origProducts[product.id].qty);
                            this.addedProducts.push(product);
                        } else if (product.qty < this.origProducts[product.id].qty) {
                            product.qty = Math.abs(this.origProducts[product.id].qty - product.qty);
                            this.removedProducts.push(product);
                        }
                    }
                }
            },
            formatProductsArray : function(productsIn) {
                var productsOut = [];
                var itemId;
                for (var i in productsIn)
                {
                    if (i != 'length' && productsIn.hasOwnProperty(i)) {
                        if (typeof(productsIn[i]['sku']) != 'undefined') {
                            itemId = productsIn[i].sku;
                        } else {
                            itemId = productsIn[i].id;
                        }
                        productsOut.push({
                            'id': itemId,
                            'name': productsIn[i].name,
                            'price': productsIn[i].price,
                            'quantity': parseInt(productsIn[i].qty, 10)
                        });
                    }
                }
                return productsOut;
            },
            cartItemAdded : function() {
                if (this.addedProducts.length == 0) {
                    return;
                }
                dataLayer.push({
                    'event': 'addToCart',
                    'ecommerce': {
                        'currencyCode' : dlCurrencyCode,
                        'add': {
                            'products': this.formatProductsArray(this.addedProducts)
                        }
                    }
                });
                this.addedProducts = [];
            },
            cartItemRemoved : function() {
                if (this.removedProducts.length == 0) {
                    return;
                }
                dataLayer.push({
                    'event': 'removeFromCart',
                    'ecommerce': {
                        'currencyCode' : dlCurrencyCode,
                        'remove': {
                            'products': this.formatProductsArray(this.removedProducts)
                        }
                    }
                });
                this.removedProducts = [];
            },
            parseAddToCartCookies : function(){
                if(getCookie(cookieAddToCart)){
                    this.addedProducts = [];
                    var addProductsList = decodeURIComponent(getCookie(cookieAddToCart));
                    this.addedProducts = JSON.parse(addProductsList);
                    delCookie(cookieAddToCart);
                    this.cartItemAdded();
                }
            },
            parseRemoveFromCartCookies : function(){
                if(getCookie(cookieRemoveFromCart)){
                    this.removedProducts = [];
                    var removeProductsList = decodeURIComponent(getCookie(cookieRemoveFromCart));
                    this.removedProducts = JSON.parse(removeProductsList);
                    delCookie(cookieRemoveFromCart);
                    this.cartItemRemoved();
                }
            }
        };

        var googleAnalyticsUniversal = new GoogleAnalyticsUniversal();
        var googleAnalyticsUniversalCart = new GoogleAnalyticsUniversalCart();

        document.observe('dom:loaded', function() {
            googleAnalyticsUniversal.updatePromotions();
            googleAnalyticsUniversal.updateImpressions();
            googleAnalyticsUniversalCart.parseAddToCartCookies();
            googleAnalyticsUniversalCart.parseRemoveFromCartCookies();
            googleAnalyticsUniversalCart.subscribeProductsUpdateInCart();
            googleAnalyticsUniversalCart.listenMinicartReload();
            dataLayer.push({'ecommerce':{'impressions':0,'promoView':0}});
        });

        function getCookie(name) {
            var cookie = " " + document.cookie;
            var search = " " + name + "=";
            var setStr = null;
            var offset = 0;
            var end = 0;
            if (cookie.length > 0) {
                offset = cookie.indexOf(search);
                if (offset != -1) {
                    offset += search.length;
                    end = cookie.indexOf(";", offset);
                    if (end == -1) {
                        end = cookie.length;
                    }
                    setStr = decodeURI(cookie.substring(offset, end));
                }
            }
            return(setStr);
        }

        function delCookie(name) {
            var date = new Date(0);
            var cookie = name + "=" + "; path=/; expires=" + date.toUTCString();
            document.cookie = cookie;
        }
    });
</script>
<script>
	require(["jquery",'recaptcha'], function($, recaptcha){
		var keys = {
			"site_key" : "6LcOVlgUAAAAAArPWuCQNqMRjO33CNmtblaUFmPP",
			"secret_key" : "6LcOVlgUAAAAAMRSCxwz4tve6U056Cg9_LMexvS3"
		};
		recaptcha.init(keys, '//www.google.com/recaptcha/api.js?render=explicit&hl=en');
	});
</script><div class="page-wrapper"><header class="page-header"><div class="header content"><span data-action="toggle-nav" class="action nav-toggle"><span>Toggle Nav</span></span>
        <a class="logo" href="https://www.sweetmarias.com/" title="Home Coffee Roasting">
        <img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/logo/stores/1/sweet-marias-homepage-logo.png"
             alt="Home Coffee Roasting"
             width="230"             height="145"        />
        <span class="logo-name-home">Home Coffee Roasting</span>
        <span class="logo-name-blog">Sweet Maria’s Coffee Library</span>
    </a>
    <div class="main-nav"><div class="nav-top"><div class="block block-search">
    <div class="block block-title"><strong>Search</strong></div>
    <div class="block block-content">
        <form class="form minisearch" id="search_mini_form" action="https://www.sweetmarias.com/catalogsearch/result/" method="get">
            <div class="field search">
                <label class="label" for="search" data-role="minisearch-label">
                    <span>Search</span>
                </label>
                <div class="control">
                    <input id="search"
                           data-mage-init='{"quickSearch":{
                                "formSelector":"#search_mini_form",
                                "url":"https://www.sweetmarias.com/search/ajax/suggest/",
                                "destinationSelector":"#search_autocomplete"}
                           }'
                           type="text"
                           name="q"
                           value=""
                           placeholder="Search"
                           class="input-text"
                           maxlength="128"
                           role="combobox"
                           aria-haspopup="false"
                           aria-autocomplete="both"
                           autocomplete="off"/>
                    <div id="search_autocomplete" class="search-autocomplete"></div>
                    <div class="nested">
    <a class="action advanced" href="https://www.sweetmarias.com/catalogsearch/advanced/" data-action="advanced-search">
        Advanced Search    </a>
</div>
                </div>
            </div>
            <div class="actions">
                <button type="submit"
                        title="Search"
                        class="action search">
                    <span>Search</span>
                </button>
            </div>
        </form>
    </div>
</div>
</div><div class="nav-center">    <div class="sections nav-sections">
                <div class="section-items nav-sections-items" data-mage-init='{"tabs":{"openedState":"active"}}'>
                                            <div class="section-item-title nav-sections-item-title" data-role="collapsible">
                    <a class="nav-sections-item-switch" data-toggle="switch" href="#store.menu">Menu</a>
                </div>
                <div class="section-item-content nav-sections-item-content" id="store.menu" data-role="content">
<nav class="navigation" data-action="navigation">
    <ul data-mage-init='{"forix/menu":{}}'>
        <li  class="level0 nav-1 first has-active level-top parent"><a href="https://www.sweetmarias.com/green-coffee.html"  class="level-top" ><span>Green Coffee</span></a><ul class="level0 submenu"><li  class="level1 nav-1-1 first parent"><a href="https://www.sweetmarias.com/green-coffee/south-america.html" ><span>South America</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-1-1-1 first"><a href="https://www.sweetmarias.com/green-coffee/south-america/colombia.html" ><span>Colombia</span></a></li><li  class="level2 nav-1-1-2"><a href="https://www.sweetmarias.com/green-coffee/south-america/brazil.html" ><span>Brazil</span></a></li><li  class="level2 nav-1-1-3 last"><a href="https://www.sweetmarias.com/green-coffee/south-america/peru.html" ><span>Peru</span></a></li></ul></li><li  class="level1 nav-1-2 parent"><a href="https://www.sweetmarias.com/green-coffee/central-america.html" ><span>Central America</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-1-2-1 first"><a href="https://www.sweetmarias.com/green-coffee/central-america/costa-rica.html" ><span>Costa Rica</span></a></li><li  class="level2 nav-1-2-2"><a href="https://www.sweetmarias.com/green-coffee/central-america/guatemala.html" ><span>Guatemala </span></a></li><li  class="level2 nav-1-2-3 last"><a href="https://www.sweetmarias.com/green-coffee/central-america/nicaragua.html" ><span>Nicaragua</span></a></li></ul></li><li  class="level1 nav-1-3 parent"><a href="https://www.sweetmarias.com/green-coffee/africa-arabia.html" ><span>Africa</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-1-3-1 first"><a href="https://www.sweetmarias.com/green-coffee/africa-arabia/burundi.html" ><span>Burundi</span></a></li><li  class="level2 nav-1-3-2"><a href="https://www.sweetmarias.com/green-coffee/africa-arabia/ethiopia.html" ><span>Ethiopia</span></a></li><li  class="level2 nav-1-3-3"><a href="https://www.sweetmarias.com/green-coffee/africa-arabia/kenya.html" ><span>Kenya</span></a></li><li  class="level2 nav-1-3-4"><a href="https://www.sweetmarias.com/green-coffee/africa-arabia/rwanda.html" ><span>Rwanda</span></a></li><li  class="level2 nav-1-3-5 last"><a href="https://www.sweetmarias.com/green-coffee/africa-arabia/congo.html" ><span>Congo</span></a></li></ul></li><li  class="level1 nav-1-4 parent"><a href="https://www.sweetmarias.com/green-coffee/indonesia-asia.html" ><span>Indonesia &amp; SE Asia</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-1-4-1 first"><a href="https://www.sweetmarias.com/green-coffee/indonesia-asia/bali.html" ><span>Bali</span></a></li><li  class="level2 nav-1-4-2"><a href="https://www.sweetmarias.com/green-coffee/indonesia-asia/sunatra.html" ><span>Sumatra</span></a></li><li  class="level2 nav-1-4-3"><a href="https://www.sweetmarias.com/green-coffee/indonesia-asia/flores.html" ><span>Flores</span></a></li><li  class="level2 nav-1-4-4"><a href="https://www.sweetmarias.com/green-coffee/indonesia-asia/java.html" ><span>Java</span></a></li><li  class="level2 nav-1-4-5 last"><a href="https://www.sweetmarias.com/green-coffee/indonesia-asia/timor.html" ><span>Timor</span></a></li></ul></li><li  class="level1 nav-1-5 has-active parent"><a href="https://www.sweetmarias.com/green-coffee/oceania.html" ><span>Oceania</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-1-5-1 first active last"><a href="https://www.sweetmarias.com/green-coffee/oceania/papua-new-guinea.html" ><span>Papua New Guinea</span></a></li></ul></li><li  class="level1 nav-1-6 last"><a href="https://www.sweetmarias.com/green-coffee/featured/new-arrivals.html" ><span>New Arrivals</span></a></li><li class="level1 parent type-menu"><p><strong><a>Shop By Type</a></strong></p>
<ul class="level1 submenu">
<ul class="level1 submenu">
<li><a href="https://www.sweetmarias.com/green-coffee/featured/new-arrivals.html">New Arrivals</a></li>
<li><a href="https://www.sweetmarias.com/green-coffee/featured/sale+coffees"Sale Coffee</a></li>
<li><a href="https://www.sweetmarias.com/green-coffee/sweet-maria-s-blends.html/">Sweet Maria&rsquo;s Blends</a></li>
<li><a href="https://www.sweetmarias.com/green-coffee/decaf.html/">Decaf</a></li>
<li><a href="/green-coffee.html?sm_status=1&amp;sm_type=77">Sample Sets</a></li>
<li><a href="/green-coffee.html?sm_flavor_profile=2058&amp;sm_status=1">Good For Espresso</a></li>
<li><a href="https://www.sweetmarias.com/green-coffee/by-type/roasted-coffee.html">Roasted Coffee</a></li>
<li><a href="https://www.sweetmarias.com/incoming-green-coffee/?___store=default">Incoming Coffees</a></li>
</ul>
</ul></li><li class="level1 parent type-image menu-link-image"><ul>
<li><a class="link-image" href="/green-coffee.html?sm_status=2"> <img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Homepage/coffee-archive.jpg" alt="Coffee Archive" width="350" height="120" /> <span>Coffee Archive</span> </a></li>
<li><a class="link-image" href="https://legacy.sweetmarias.com/library/category/origins/origin-pages"> <img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Homepage/origin-pages.jpg" alt="Origin Pages" width="350" height="120" /> <span>Origin Pages</span> </a></li>
<li><a class="link-image" href="https://legacy.sweetmarias.com/library/how-to-roast-your-own-coffee"> <img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Homepage/home-roasting-basics.jpg" alt="Coffee Archive" width="350" height="120" /> <span>Home Roasting Basics</span> </a></li>
</ul></li></ul></li><li  class="level0 nav-2 level-top parent"><a href="https://www.sweetmarias.com/roasting.html"  class="level-top" ><span>Roasting</span></a><ul class="level0 submenu"><li  class="level1 nav-2-1 first parent"><a href="https://www.sweetmarias.com/roasting/air-roasters.html" ><span>Air Roasters</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-2-1-1 first"><a href="https://www.sweetmarias.com/roasting/air-roasters/freshroast.html" ><span>Freshroast</span></a></li><li  class="level2 nav-2-1-2"><a href="https://www.sweetmarias.com/roasting/air-roasters/nesco.html" ><span>Nesco</span></a></li><li  class="level2 nav-2-1-3 last"><a href="https://www.sweetmarias.com/roasting/air-roasters/hot-air-poppers.html" ><span>Hot Air Poppers</span></a></li></ul></li><li  class="level1 nav-2-2 parent"><a href="https://www.sweetmarias.com/roasting/drum-roasters.html" ><span>Drum Roasters</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-2-2-1 first"><a href="https://www.sweetmarias.com/roasting/drum-roasters/hottop.html" ><span>Hottop</span></a></li><li  class="level2 nav-2-2-2"><a href="https://www.sweetmarias.com/roasting/drum-roasters/behmor.html" ><span>Behmor</span></a></li><li  class="level2 nav-2-2-3 last"><a href="https://www.sweetmarias.com/roasting/drum-roasters/gene-cafe.html" ><span>Gene Cafe</span></a></li></ul></li><li  class="level1 nav-2-3 parent"><a href="https://www.sweetmarias.com/roasting/stovetop.html" ><span>Stovetop Roasters</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-2-3-1 first last"><a href="https://www.sweetmarias.com/roasting/stovetop/stovepop-stainless-steel.html" ><span>StovePop-Stainless Steel</span></a></li></ul></li><li  class="level1 nav-2-4 last parent"><a href="https://www.sweetmarias.com/roasting/roasting-accessories.html" ><span>Roasting Accessories</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-2-4-1 first last"><a href="https://www.sweetmarias.com/roasting/roasting-accessories/coffee-storage.html" ><span>Coffee Storage</span></a></li></ul></li><li class="level1 parent type-image menu-link-image"><ul>
<li><a class="link-image" href="/roasting/starter-kits.html"> <img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Homepage/sweet-marias-starter-kits.jpg" alt="Roasting Kits" width="350" height="230" /> <span>Starter Kits</span> </a></li>
</ul></li></ul></li><li  class="level0 nav-3 level-top parent"><a href="https://www.sweetmarias.com/brewing.html"  class="level-top" ><span>Brewing</span></a><ul class="level0 submenu full-width"><li  class="level1 nav-3-1 first parent"><a href="https://www.sweetmarias.com/brewing/brewers.html" ><span>Brewers</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-3-1-1 first"><a href="https://www.sweetmarias.com/brewing/brewers/electric.html" ><span>Electric</span></a></li><li  class="level2 nav-3-1-2"><a href="https://www.sweetmarias.com/brewing/brewers/pour-over.html" ><span>Pour Over</span></a></li><li  class="level2 nav-3-1-3"><a href="https://www.sweetmarias.com/brewing/brewers/vacuum.html" ><span>Vacuum</span></a></li><li  class="level2 nav-3-1-4 last"><a href="https://www.sweetmarias.com/brewing/brewers/coffee-presses.html" ><span>Coffee Presses</span></a></li></ul></li><li  class="level1 nav-3-2 parent"><a href="https://www.sweetmarias.com/brewing/coffee-drinkware.html" ><span>Coffee Drinkware</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-3-2-1 first"><a href="https://www.sweetmarias.com/brewing/coffee-drinkware/coffee-cups.html" ><span>Coffee Cups</span></a></li><li  class="level2 nav-3-2-2 last"><a href="https://www.sweetmarias.com/brewing/coffee-drinkware/thermal-cups-bottles.html" ><span>Thermal Cups &amp; Bottles</span></a></li></ul></li><li  class="level1 nav-3-3 parent"><a href="https://www.sweetmarias.com/brewing/espresso.html" ><span>Espresso</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-3-3-1 first"><a href="https://www.sweetmarias.com/brewing/espresso/espresso-machines.html" ><span>Espresso Machines</span></a></li><li  class="level2 nav-3-3-2 last"><a href="https://www.sweetmarias.com/brewing/espresso/espresso-accessories.html" ><span>Espresso Accessories</span></a></li></ul></li><li  class="level1 nav-3-4 parent"><a href="https://www.sweetmarias.com/brewing/grinders.html" ><span>Grinders</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-3-4-1 first"><a href="https://www.sweetmarias.com/brewing/grinders/electric.html" ><span>Electric</span></a></li><li  class="level2 nav-3-4-2 last"><a href="https://www.sweetmarias.com/brewing/grinders/manual.html" ><span>Manual</span></a></li></ul></li><li  class="level1 nav-3-5 last parent"><a href="https://www.sweetmarias.com/brewing/accessories.html" ><span>Accessories</span></a><ul class="level1 submenu full-width"><li  class="level2 nav-3-5-1 first"><a href="https://www.sweetmarias.com/brewing/accessories/cleaning-supplies.html" ><span>Cleaning Supplies</span></a></li><li  class="level2 nav-3-5-2"><a href="https://www.sweetmarias.com/brewing/accessories/filters.html" ><span>Filters</span></a></li><li  class="level2 nav-3-5-3 last"><a href="https://www.sweetmarias.com/brewing/accessories/kettles.html" ><span>Kettles</span></a></li></ul></li></ul></li><li  class="level0 nav-4 level-top"><a href="https://www.sweetmarias.com/extras.html"  class="level-top" ><span>Extras</span></a><ul class="level0 submenu"><li class="level1 parent type-image menu-link-image"><ul>
<li><a class="link-image" href="https://www.sweetmarias.com/extras/this-and-that.html/"> <img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Homepage/this-and-that.jpg" alt="This And That" width="350" height="230" /><span>This And That</span></a> <a class="abutton" href="https://www.sweetmarias.com/extras/this-and-that.html/">SHOP NOW</a></li>
<li><a class="link-image" href="https://www.sweetmarias.com/extras/apparel.html/"> <img title="Apparel" src="https://legacy.sweetmarias.com/sweet-blog/wp-content/uploads/2018/09/sweet-marias-apparel.jpg" alt="Apparel" width="350" height="230" /><span>Apparel</span></a>
<p><a class="abutton" href="https://www.sweetmarias.com/extras/apparel.html/">SHOP NOW</a></p>
</li>
<li><a class="link-image" href="https://www.sweetmarias.com/extras/books.html/"> <img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Homepage/books.jpg" alt="Books" width="350" height="230" /><span>Books</span> </a>
<p><a class="abutton" href="https://www.sweetmarias.com/extras/books.html/">SHOP NOW</a></p>
</li>
<li><a class="link-image" href="https://www.sweetmarias.com/gift-card.html/"><img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/gift-certificates.jpg" alt="Gift Certificates" width="350" height="230" /><span>Gift Certificates</span> </a>
<p><a class="abutton" href="https://www.sweetmarias.com/gift-card.html/">SHOP NOW</a></p>
</li>
</ul></li></ul></li><li  class="level0 nav-5 level-top"><a href="https://legacy.sweetmarias.com/library/"  class="level-top" ><span>Resources</span></a><ul class="level0 submenu"><li class="level1 parent type-image menu-link-image"><ul>
<li><a class="link-image" href="https://legacy.sweetmarias.com/library/how-to-roast-your-own-coffee"> <img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Homepage/home-coffee-roasting-basics.jpg" alt="Roasting Basics" width="350" height="230" /> <span>Roasting Basics</span> </a>
<p>Home roasting can be easy and simple. Let's get started.</p>
<!--<a class="link-type3" href="#">Roasting Basics Guide</a> <a class="link-type3" href="#">Home Roasting Forum</a></li>--></li>
<li><a class="link-image" href="https://legacy.sweetmarias.com/library/category/video/"> <img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Homepage/media.jpg" alt="Media" width="350" height="230" /> <span>Media</span> </a>
<p>Videos, photos, and podcasts about coffee and roasting.</p>
<!--<a class="link-type3" href="#">Product Videos</a> <a class="link-type3" href="#">Sweet Maria&rsquo;s Podcast</a>--></li>
<li><a class="link-image" href="https://www.sweetmarias.com/glossary/"> <img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Homepage/coffee-glossary-small.jpg" alt="Coffee Glossary" width="350" height="230" /><span>Coffee Glossary</span></a>
<p>Every coffee word we could think of...defined.</p>
<!--<a class="link-type3" href="#">Link 1</a> <a class="link-type3" href="#">Link 2</a>--></li>
<li><a class="link-image" href="https://legacy.sweetmarias.com/library/category/origins/origin-pages/"> <img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Homepage/origin-pages_1.jpg" alt="Origin Pages" width="350" height="230" /><span>Origin Pages</span></a>
<p>Tom's insights on every country we import coffee from</p>
<!--<a class="link-type3" href="#">Link 1</a> <a class="link-type3" href="#">Link 2</a>--></li>
<li class="bottom-area">
<p><a href="https://legacy.sweetmarias.com/library/" target="_self">Check out our coffee library!</a></p>
</li>
</ul></li></ul></li><li  class="level0 nav-6 last level-top"><a href="https://www.sweetmarias.com/sale-items-miscellaneous-brewing-kits-brewing-equipment-brewers-aeropress.html"  class="level-top" ><span>Sale </span></a></li>            </ul>
</nav>
</div>
                                            <div class="section-item-title nav-sections-item-title" data-role="collapsible">
                    <a class="nav-sections-item-switch" data-toggle="switch" href="#store.links">Account</a>
                </div>
                <div class="section-item-content nav-sections-item-content" id="store.links" data-role="content"><!-- Account links --></div>
                                    </div>
    </div>
</div></div><div data-block="minicart" class="minicart-wrapper">
    <a class="action showcart" href="https://www.sweetmarias.com/checkout/cart/"
       data-bind="scope: 'minicart_content'" data-trigger="minicart" title="My Cart">
        <span class="text">Cart</span>
        <span class="counter qty empty"
              data-bind="css: { empty: !!getCartParam('summary_count') == false }, blockLoader: isLoading">
            <span class="counter-number" data-bind="attr: {'data-number': getCartParam('summary_count')}">
                <!-- ko if: getCartParam('summary_count') -->
                    <!-- ko text: getCartParam('summary_count') --><!-- /ko -->
                <!-- /ko -->
                <!-- ko ifnot: getCartParam('summary_count') -->
                    <!-- ko text: 0 --><!-- /ko -->
                <!-- /ko -->
            </span>
            <span class="counter-label">
                <!-- ko if: getCartParam('summary_count') -->
                    <!-- ko text: getCartParam('summary_count') --><!-- /ko -->
                    <!-- ko i18n: 'items' --><!-- /ko -->
                <!-- /ko -->
                <!-- ko ifnot: getCartParam('summary_count') -->
                    <!-- ko text: 0 --><!-- /ko -->
                    <!-- ko i18n: 'item' --><!-- /ko -->
                <!-- /ko -->
            </span>

            <div class="mini-cart-loading-container" style="display: none;">
                <div class="loader">
                    <img class="mini-cart-loading-img" src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/static/version1553065162/frontend/Forix/sweetmarias/en_US/images/loader-1.gif" alt="Loading...">
                </div>
            </div>
        </span>
    </a>
            <div class="block block-minicart empty swipe-minicart"
             data-role="dropdownDialog" data-bind="mageInit: {
            'forix/pushdata':{
                'container': '.minicart-wrapper',
                'toggleBtn': '[data-trigger=minicart]',
                'swipeArea': '.swipe-minicart',
                'pushCloseCls':'push-close',
                'closeBtnCls':'toggle-cart-close',
                'noEffect':'#mini-login',
                'clsPush':'minicart',
                'responsive': true,
                'swiped': 'right'
            }}">
            <div id="minicart-content-wrapper" data-bind="scope: 'minicart_content'">
                <!-- ko template: getTemplate() --><!-- /ko -->
                            </div>
                    </div>
        <script>
        window.checkout = {"shoppingCartUrl":"https:\/\/www.sweetmarias.com\/checkout\/cart\/","checkoutUrl":"https:\/\/www.sweetmarias.com\/checkout\/","updateItemQtyUrl":"https:\/\/www.sweetmarias.com\/checkout\/sidebar\/updateItemQty\/","removeItemUrl":"https:\/\/www.sweetmarias.com\/checkout\/sidebar\/removeItem\/","imageTemplate":"Magento_Catalog\/product\/image_with_borders","baseUrl":"https:\/\/www.sweetmarias.com\/","minicartMaxItemsVisible":2,"websiteId":"1","maxItemsToDisplay":10,"customerLoginUrl":"https:\/\/www.sweetmarias.com\/customer\/account\/login\/","isRedirectRequired":false,"autocomplete":"off","captcha":{"user_login":{"isCaseSensitive":false,"imageHeight":50,"imageSrc":"","refreshUrl":"https:\/\/www.sweetmarias.com\/captcha\/refresh\/","isRequired":false},"guest_checkout":{"isCaseSensitive":false,"imageHeight":50,"imageSrc":"","refreshUrl":"https:\/\/www.sweetmarias.com\/captcha\/refresh\/","isRequired":false}}};
    </script>
    <script type="text/x-magento-init">
    {
        "[data-block='minicart']": {
            "Magento_Ui/js/core/app": {"components":{"minicart_content":{"children":{"subtotal.container":{"children":{"subtotal":{"children":{"subtotal.totals":{"config":{"display_cart_subtotal_incl_tax":0,"display_cart_subtotal_excl_tax":1,"template":"Magento_Tax\/checkout\/minicart\/subtotal\/totals"},"component":"Magento_Tax\/js\/view\/checkout\/minicart\/subtotal\/totals","children":{"subtotal.totals.msrp":{"component":"Magento_Msrp\/js\/view\/checkout\/minicart\/subtotal\/totals","config":{"displayArea":"minicart-subtotal-hidden","template":"Magento_Msrp\/checkout\/minicart\/subtotal\/totals"}}}}},"component":"uiComponent","config":{"template":"Magento_Checkout\/minicart\/subtotal"}}},"component":"uiComponent","config":{"displayArea":"subtotalContainer"}},"item.renderer":{"component":"uiComponent","config":{"displayArea":"defaultRenderer","template":"Magento_Checkout\/minicart\/item\/default"},"children":{"item.image":{"component":"Magento_Catalog\/js\/view\/image","config":{"template":"Magento_Catalog\/product\/image","displayArea":"itemImage"}},"checkout.cart.item.price.sidebar":{"component":"uiComponent","config":{"template":"Magento_Checkout\/minicart\/item\/price","displayArea":"priceSidebar"}}}},"extra_info":{"component":"uiComponent","config":{"displayArea":"extraInfo"}},"promotion":{"component":"uiComponent","config":{"displayArea":"promotion"}}},"config":{"itemRenderer":{"default":"defaultRenderer","simple":"defaultRenderer","virtual":"defaultRenderer"},"template":"Magento_Checkout\/minicart\/content"},"component":"Magento_Checkout\/js\/view\/minicart"}},"types":[]}        },
        "*": {
            "Magento_Ui/js/block-loader": "https://sweet-marias-live-sweetmarias.netdna-ssl.com/static/version1553065162/frontend/Forix/sweetmarias/en_US/images/loader-1.gif"
        }
    }
    </script>
</div>


<div id="mini-login" class="block minilogin-wrapper"  data-block="mini_login" data-bind="scope: 'mini_login_content'">
        <input name="form_key" type="hidden" value="2hKQSKcftqnlC4Bx" />    <a class="action showlogin" data-trigger="authentication" href="javascript:void(0)" title="Login"><span data-bind="i18n: 'Log In'"></span></a>
    <!-- ko template: getTemplate() --><!-- /ko -->
            <a class="register" href="https://www.sweetmarias.com/customer/account/create/" title="Register"><span data-bind="i18n: 'Register'"></span></a>
    
    <script>
        window.customerData = [];
    </script>
    <script type="text/x-magento-init">
    {
        "#mini-login": {
            "Magento_Ui/js/core/app": {"components":{"mini_login_content":{"component":"Forix_Minilogin\/js\/view\/minilogin","config":{"template":"Forix_Minilogin\/login"},"children":{"messages":{"component":"Magento_Ui\/js\/view\/messages","displayArea":"messages"}}}}}        },
        "*": {
            "Magento_Ui/js/block-loader": "https://sweet-marias-live-sweetmarias.netdna-ssl.com/static/version1553065162/frontend/Forix/sweetmarias/en_US/images/loader-1.gif"
        }
    }
    </script>
    </div>



</div></header><div class="breadcrumbs" data-block="breadcrumbs">
    <ul class="items">
                    <li class="item home ">
                            <a href="https://www.sweetmarias.com/" title="Go to Home Page">
                    Home                </a>
                        </li>
                    <li class="item category4 ">
                            <a href="https://www.sweetmarias.com/green-coffee.html" title="">
                    Green Coffee                </a>
                        </li>
                    <li class="item category57 ">
                            <a href="https://www.sweetmarias.com/green-coffee/oceania.html" title="">
                    Oceania                </a>
                        </li>
                    <li class="item category231 ">
                            <a href="https://www.sweetmarias.com/green-coffee/oceania/papua-new-guinea.html" title="">
                    Papua New Guinea                </a>
                        </li>
                    <li class="item product last">
                            <strong>Papua New Guinea Honey Nebilyer Estate</strong>
                        </li>
            </ul>
</div>
<div class="page messages"><div data-placeholder="messages"></div>
<div data-bind="scope: 'messages'">
    <div data-bind="foreach: { data: cookieMessages, as: 'message' }" class="messages">
        <div data-bind="attr: {
            class: 'message-' + message.type + ' ' + message.type + ' message',
            'data-ui-id': 'message-' + message.type
        }">
            <div data-bind="html: message.text"></div>
        </div>
    </div>
    <div data-bind="foreach: { data: messages().messages, as: 'message' }" class="messages">
        <div data-bind="attr: {
            class: 'message-' + message.type + ' ' + message.type + ' message',
            'data-ui-id': 'message-' + message.type
        }">
            <div data-bind="html: message.text"></div>
        </div>
    </div>
</div>
<script type="text/x-magento-init">
    {
        "*": {
            "Magento_Ui/js/core/app": {
                "components": {
                        "messages": {
                            "component": "Magento_Theme/js/view/messages"
                        }
                    }
                }
            }
    }
</script>
</div><main id="maincontent" class="page-main"><a id="contentarea" tabindex="-1"></a>
<div class="columns"><div class="column main"><div class="product media">        <span class="videos-list"></span>
    <div class="gallery-placeholder _block-content-loading green-coffee-gallery" data-gallery-role="gallery-placeholder">
        <div data-role="loader" class="loading-mask">
            <div class="loader">
                <img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/static/version1553065162/frontend/Forix/sweetmarias/en_US/images/loader-1.gif"
                     alt="Loading...">
            </div>
        </div>
    </div>
    <!--Fix for jumping content. Loader must be the same size as gallery.-->
            <script>
            require(["jquery","domReady!"], function ($) {
                $("body").addClass("product-full-gallery");
            });

            var config = {
                    "width": 1338,
                    "thumbheight": 178,
                    "navtype": "thumbs",
                    "height": 576                },
                thumbBarHeight = 0,
                loader = document.querySelectorAll('[data-gallery-role="gallery-placeholder"] [data-role="loader"]')[0];

            if (config.navtype === 'horizontal') {
                thumbBarHeight = config.thumbheight;
            }

            loader.style.paddingBottom = ( config.height / config.width * 100) + "";
        </script>

        <script type="text/x-magento-init">
        {
            "[data-gallery-role=gallery-placeholder]": {
                "mage/gallery/gallery-ext": {
                    "mixins":["magnifier/magnify"],
                    "magnifierOpts": {"fullscreenzoom":"4","top":"","left":"","width":"","height":"","eventType":"hover","enabled":"false"},
                    "data": [{"thumb":"https:\/\/sweet-marias-live-sweetmarias.netdna-ssl.com\/media\/catalog\/product\/cache\/thumbnail\/178x178\/beff4985b56e3afdbeabfc89641a4582\/p\/a\/papua-new-guinea-sweetmarias-2.jpg","img":"https:\/\/sweet-marias-live-sweetmarias.netdna-ssl.com\/media\/catalog\/product\/cache\/image\/1338x576\/beff4985b56e3afdbeabfc89641a4582\/p\/a\/papua-new-guinea-sweetmarias-2.jpg","full":"https:\/\/sweet-marias-live-sweetmarias.netdna-ssl.com\/media\/catalog\/product\/cache\/image\/beff4985b56e3afdbeabfc89641a4582\/p\/a\/papua-new-guinea-sweetmarias-2.jpg","caption":null,"position":"6","isMain":true},{"thumb":"https:\/\/sweet-marias-live-sweetmarias.netdna-ssl.com\/media\/catalog\/product\/cache\/thumbnail\/178x178\/beff4985b56e3afdbeabfc89641a4582\/p\/a\/papua-new-guinea-sweetmarias-4.jpg","img":"https:\/\/sweet-marias-live-sweetmarias.netdna-ssl.com\/media\/catalog\/product\/cache\/image\/1338x576\/beff4985b56e3afdbeabfc89641a4582\/p\/a\/papua-new-guinea-sweetmarias-4.jpg","full":"https:\/\/sweet-marias-live-sweetmarias.netdna-ssl.com\/media\/catalog\/product\/cache\/image\/beff4985b56e3afdbeabfc89641a4582\/p\/a\/papua-new-guinea-sweetmarias-4.jpg","caption":null,"position":"7","isMain":false},{"thumb":"https:\/\/sweet-marias-live-sweetmarias.netdna-ssl.com\/media\/catalog\/product\/cache\/thumbnail\/178x178\/beff4985b56e3afdbeabfc89641a4582\/p\/a\/papua-new-guinea-sweetmarias-14.jpg","img":"https:\/\/sweet-marias-live-sweetmarias.netdna-ssl.com\/media\/catalog\/product\/cache\/image\/1338x576\/beff4985b56e3afdbeabfc89641a4582\/p\/a\/papua-new-guinea-sweetmarias-14.jpg","full":"https:\/\/sweet-marias-live-sweetmarias.netdna-ssl.com\/media\/catalog\/product\/cache\/image\/beff4985b56e3afdbeabfc89641a4582\/p\/a\/papua-new-guinea-sweetmarias-14.jpg","caption":null,"position":"8","isMain":false}],
                    "options": {
                        "nav": "thumbs",
                                                    "loop": true,
                                                                            "keyboard": false,
                                                                            "arrows": false,
                                                                            "allowfullscreen": true,
                                                                            "showCaption": true,
                                                "width": "1338",
                                                    "height": 576,
                                                "thumbwidth": "178",
                                                    "thumbheight": 178,
                                                                            "transitionduration": 500,
                                                "transition": "crossfade",
                                                    "navarrows": true,
                                                "navtype": "thumbs",
                        "navdir": "horizontal"
                    },
                    "fullscreen": {
                        "nav": "thumbs",
                                                    "loop": true,
                                                                            "keyboard": true,
                                                                            "showCaption": true,
                                                "navdir": "horizontal",
                                                "navtype": "slides",
                                                    "arrows": true,
                                                                            "showCaption": true,
                                                                            "transitionduration": 500,
                                                "transition": "slides",
                        "navwidth": "100%"
                    },
                    "breakpoints": {"desktop":{"conditions":{"min-width":"1024px"},"options":{"options":{"navwidth":"53.21%"}}},"tablet":{"conditions":{"max-width":"1023px","min-width":"768px"},"options":{"options":{"navwidth":"100%","thumbwidth":"192","thumbheight":"192"}}},"mobile":{"conditions":{"max-width":"767px"},"options":{"options":{"nav":"dots","navwidth":"100%"}}}}                }
            }
        }
        </script>
    <!-- Product not Green Coffee-->
    
    <script type="text/javascript">
        require(['jquery', 'mage/gallery/gallery'], function($, gallery){
            $('.product.media [data-gallery-role=gallery-placeholder]').on('gallery:loaded', function () {
                //if($(window).width() < 768){
//                var  api;
                $(this).on('fotorama:ready', function(){
                    var api = $(this).data('gallery');

                    if(!$(".fotorama__nav-wrap").hasClass("hasButtonArrow")){
                        $(".fotorama__nav-wrap").addClass("hasButtonArrow");

                        $(".fotorama__thumb__arr.fotorama__thumb__arr--left").clone().addClass("gallery-prev").prependTo($(".fotorama__nav-wrap"));
                        $(".fotorama__thumb__arr.fotorama__thumb__arr--right").clone().addClass("gallery-next").appendTo($(".fotorama__nav-wrap"));

                        $(".gallery-prev").click(function(){
                            api.prev();
                        });

                        $(".gallery-next").click(function(){
                            api.next();
                        });

                        $(".fotorama__stage__shaft").append("<div class='zoom-fullscreen' title='Click To Zoom Image'></div>");
                    }

//                    if(api.fotorama.options.navdir == 'vertical'){
//                        api.fotorama.options.navdir = 'horizontal';
//                        api.fotorama.resize();
//                    }

//                    $(".product.media .video-thumb-icon").click(function(){
//                        if(!$(this).hasClass("playingVideo")){
//                            $(this).addClass("playingVideo");
//                            $(".fotorama-video-container").trigger("click");
//
//                            console.log("PlayingVideo");
//                        }
//                        console.log("Click Thumb Video");
//                        return false;
//                    });
                });
                //}

            });
        });
    </script>

<script type="text/x-magento-init">
    {
        "[data-gallery-role=gallery-placeholder]": {
            "Magento_ProductVideo/js/fotorama-add-video-events": {
                "videoData": [{"mediaType":"image","videoUrl":null,"isBase":true},{"mediaType":"image","videoUrl":null,"isBase":false},{"mediaType":"image","videoUrl":null,"isBase":false}],
                "videoSettings": [{"playIfBase":"1","showRelated":"0","videoAutoRestart":"1"}],
                "optionsVideoData": []            }
        }
    }
</script>
</div><div class="product-main-container"><div class="product-info-main"><div class="page-title-wrapper product">
    <h1 class="page-title"
                >
        <span class="base" data-ui-id="page-title-wrapper" itemprop="name">Papua New Guinea Honey Nebilyer Estate</span>    </h1>
    </div>

<div class="product attribute overview">
        <div class="value" itemprop="description"><p>Honey process seems to bring out fruited notes like cranberry, raisin, red grape, and underscored by molasses-like sweetness. This PNG boasts body, and with mild acidity, is great espresso too. City+ to Full City+. Good for espresso.</p></div>
</div>


<div class="product-add-form">
    <form action="https://www.sweetmarias.com/checkout/cart/add/uenc/aHR0cHM6Ly93d3cuc3dlZXRtYXJpYXMuY29tL3BhcHVhLW5ldy1ndWluZWEtaG9uZXktcHJvY2Vzcy1uZWJpbHllci1lc3RhdGUuaHRtbA,,/product/14051/" method="post"
          id="product_addtocart_form">
        <input type="hidden" name="product" value="14051" />
        <input type="hidden" name="selected_configurable_option" value="" />
        <input type="hidden" name="related_product" id="related-products-field" value="" />
        <input name="form_key" type="hidden" value="2hKQSKcftqnlC4Bx" />                                    
                    <div class="product-options-wrapper" id="product-options-wrapper" data-hasrequired="* Required Fields">
    <div class="fieldset" tabindex="0">
        <div class="product-info-price"><div class="price-box price-final_price" data-role="priceBox" data-product-id="14051">
    

<span class="price-container price-final_price tax weee"
         itemprop="offers" itemscope itemtype="http://schema.org/Offer">
        <span  id="product-price-14051"                data-price-amount="12.07"
        data-price-type="finalPrice"
        class="price-wrapper "
         itemprop="price">
        <span class="price">$12.07</span>    </span>
                <meta itemprop="priceCurrency" content="USD" />
    </span>
</div><div class="product-info-stock-sku">
            <div class="stock available" title="Availability">
            <span>In stock</span>
        </div>
    </div></div>
            <div class="field configurable required">
            <label class="label" for="attribute190">
                <span>Weight</span>
            </label>
            <div class="control">
                <select name="super_attribute[190]"
                        data-selector="super_attribute[190]"
                        data-validate="{required:true}"
                        id="attribute190"
                        class="super-attribute-select">
                    <option value="">Choose an Option...</option>
                </select>
            </div>
        </div>
        <script type="text/x-magento-init">
        {
            "#product_addtocart_form": {
                "Magento_ConfigurableProduct/js/configurableExt": {
                    "spConfig": {"attributes":{"190":{"id":"190","code":"sm_weight","label":"Weight","options":[{"id":"664","label":"50 LB","products":[]},{"id":"665","label":"100 LB","products":[]},{"id":"110","label":" 1 LB","products":["14233"]},{"id":"660","label":"2 LB","products":["14234"]},{"id":"661","label":"5 LB","products":["14235"]},{"id":"662","label":"10 LB","products":["14236"]},{"id":"663","label":"20 LB","products":["14237"]}],"position":"0"}},"template":"$<data.price>","optionPrices":{"14238":{"oldPrice":{"amount":"232.5"},"basePrice":{"amount":"232.5"},"finalPrice":{"amount":"232.5"}},"14239":{"oldPrice":{"amount":"450"},"basePrice":{"amount":"450"},"finalPrice":{"amount":"450"}},"14233":{"oldPrice":{"amount":"6.35"},"basePrice":{"amount":"6.35"},"finalPrice":{"amount":"6.35"}},"14234":{"oldPrice":{"amount":"12.07"},"basePrice":{"amount":"12.07"},"finalPrice":{"amount":"12.07"}},"14235":{"oldPrice":{"amount":"27.63"},"basePrice":{"amount":"27.63"},"finalPrice":{"amount":"27.63"}},"14236":{"oldPrice":{"amount":"52.71"},"basePrice":{"amount":"52.71"},"finalPrice":{"amount":"52.71"}},"14237":{"oldPrice":{"amount":"100.79"},"basePrice":{"amount":"100.79"},"finalPrice":{"amount":"100.79"}}},"prices":{"oldPrice":{"amount":"12.07"},"basePrice":{"amount":"12.07"},"finalPrice":{"amount":"12.07"}},"productId":"14051","chooseText":"Choose an Option...","images":[],"index":{"14233":{"190":"110"},"14234":{"190":"660"},"14235":{"190":"661"},"14236":{"190":"662"},"14237":{"190":"663"}}},
                    "onlyMainImg": true                }
            }
        }
    </script>

    <script type="text/x-magento-init">
                    {
                        ".product-options-wrapper": {
                                    "Amasty_Xnotif/js/amnotification": {
                                        "xnotif": {"110":{"is_in_stock":false,"custom_status":"Out of Stock","product_id":"14233"},"660":{"is_in_stock":false,"custom_status":"Out of Stock","product_id":"14234"},"661":{"is_in_stock":false,"custom_status":"Out of Stock","product_id":"14235"},"662":{"is_in_stock":false,"custom_status":"Out of Stock","product_id":"14236"},"663":{"is_in_stock":false,"custom_status":"Out of Stock","product_id":"14237"},"changeConfigurableStatus":true}
                                    }
                         }
                    }
                   </script>
<script>
require([
    "jquery",
    "jquery/ui"
], function($){

//<![CDATA[
    $.extend(true, $, {
        calendarConfig: {
            dayNames: ["Sunday","Monday","Tuesday","Wednesday","Thursday","Friday","Saturday"],
            dayNamesMin: ["Sun","Mon","Tue","Wed","Thu","Fri","Sat"],
            monthNames: ["January","February","March","April","May","June","July","August","September","October","November","December"],
            monthNamesShort: ["Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"],
            infoTitle: "About the calendar",
            firstDay: 0,
            closeText: "Close",
            currentText: "Go Today",
            prevText: "Previous",
            nextText: "Next",
            weekHeader: "WK",
            timeText: "Time",
            hourText: "Hour",
            minuteText: "Minute",
            dateFormat: $.datepicker.RFC_2822,
            showOn: "button",
            showAnim: "",
            changeMonth: true,
            changeYear: true,
            buttonImageOnly: null,
            buttonImage: null,
            showButtonPanel: true,
            showWeek: true,
            timeFormat: '',
            showTime: false,
            showHour: false,
            showMinute: false
        }
    });

    enUS = {"m":{"wide":["January","February","March","April","May","June","July","August","September","October","November","December"],"abbr":["Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"]}}; // en_US locale reference
//]]>

});
</script>

<div class="product-options-bottom">
    <div class="box-tocart">
    <div class="fieldset">
                <div class="field qty">
            <label class="label" for="qty"><span>Qty</span></label>
            <div class="control">
                <input type="number"
                       name="qty"
                       id="qty"
                       maxlength="12"
                       value="1"
                       title="Qty" class="input-text qty"
                       data-validate="{&quot;required-number&quot;:true,&quot;validate-item-quantity&quot;:{&quot;minAllowed&quot;:1}}"
                       />
            </div>
        </div>
                <div class="actions">
            <button type="submit"
                    title="Add to Cart"
                    class="action primary tocart"
                    id="product-addtocart-button">
                <span>Add to Cart</span>
            </button>
                    </div>
    </div>
</div>
<script>
    require([
        'jquery',
        'mage/mage',
        'Magento_Catalog/product/view/validation',
        'Magento_Catalog/js/catalog-add-to-cart-ext'
    ], function ($) {
        'use strict';

        $('#product_addtocart_form').mage('validation', {
            radioCheckboxClosest: '.nested',
            submitHandler: function (form) {
                var widget = $(form).catalogAddToCartExt({
                    bindSubmit: false
                });

                widget.catalogAddToCartExt('submitForm', $(form));

                return false;
            }
        });
    });
</script>
<div class="product-social-links"><div class="product-addto-links" data-role="add-to-links">
        <a href="#"
       class="action towishlist"
       data-post='{"action":"https:\/\/www.sweetmarias.com\/wishlist\/index\/add\/","data":{"product":"14051","uenc":"aHR0cHM6Ly93d3cuc3dlZXRtYXJpYXMuY29tL3BhcHVhLW5ldy1ndWluZWEtaG9uZXktcHJvY2Vzcy1uZWJpbHllci1lc3RhdGUuaHRtbA,,"}}'
       title="Add To Wishlist" 
       data-action="add-to-wishlist"><span>Wish List</span></a>
<script type="text/x-magento-init">
    {
        "body": {
            "addToWishlist": {"productType":"configurable","giftcardInfo":"[id^=giftcard]"}        }
    }
</script>
<div class="action shareto">
    <span>Share</span>
    <div class="share-container">
        <div class="addthis_inline_share_toolbox"></div>
            <div class="at-style-responsive at-resp-share-element at-share-btn-elements">
        <a role="button" href="https://www.sweetmarias.com/sendfriend/product/send/id/14051/cat_id/231/"
           tabindex="1" class="at-icon-wrapper at-share-btn at-svc-email action mailto at-share-btn friend fancybox fancybox.iframe' : '' ?>" style="background-color: rgb(132, 132, 132); border-radius: 8px;">
            <span class="at4-visually-hidden">Share to Email</span>
                    <span class="at-icon-wrapper" style="line-height: 32px; height: 32px; width: 32px;">
                        <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 32 32" style="fill: rgb(255, 255, 255); width: 32px; height: 32px;" class="at-icon at-icon-email"><g><g fill-rule="evenodd"></g>
                                <path d="M27 22.757c0 1.24-.988 2.243-2.19 2.243H7.19C5.98 25 5 23.994 5 22.757V13.67c0-.556.39-.773.855-.496l8.78 5.238c.782.467 1.95.467 2.73 0l8.78-5.238c.472-.28.855-.063.855.495v9.087z"></path><path d="M27 9.243C27 8.006 26.02 7 24.81 7H7.19C5.988 7 5 8.004 5 9.243v.465c0 .554.385 1.232.857 1.514l9.61 5.733c.267.16.8.16 1.067 0l9.61-5.733c.473-.283.856-.96.856-1.514v-.465z"></path></g>
                        </svg>
                </span>
            <span class="at-label" style="font-size: 11.4px; line-height: 32px; height: 32px; color: rgb(255, 255, 255);">Email</span></a>
    </div>
    <script type="text/x-magento-init">
    {
        ".fancybox": {
            "forix/fancybox": {
                            
            }
        }
    }
    </script>
    </div>
</div>
</div>
</div></div>
    </div>
</div>
                                </form>
</div>
<script>
    require([
        'jquery',
        'priceBox'
    ], function($){
        var dataPriceBoxSelector = '[data-role=priceBox]',
            dataProductIdSelector = '[data-product-id=14051]',
            priceBoxes = $(dataPriceBoxSelector + dataProductIdSelector);

        priceBoxes = priceBoxes.filter(function(index, elem){
            return !$(elem).find('.price-from').length;
        });

        priceBoxes.priceBox({'priceConfig': {"productId":"14051","priceFormat":{"pattern":"$","precision":2,"requiredPrecision":2,"decimalSymbol":".","groupSymbol":",","groupLength":3,"integerRequired":1},"prices":{"oldPrice":{"amount":12.07,"adjustments":[]},"basePrice":{"amount":12.07,"adjustments":[]},"finalPrice":{"amount":12.07,"adjustments":[]}},"idSuffix":"_clone","tierPrices":[],"calculationAlgorithm":"TOTAL_BASE_CALCULATION"}});
    });
</script>
    <div class="product-custom-attribute">
                    <div class="score-origin">
                                    <div class="total-score">
                        <h5 class="score-value">86.6</h5>
                        <label class="score-label">Total Score</label>
                    </div>
                                            </div>
                
        
                    </div>
</div></div>    <div class="product info detailed">
                <div class="product data items" data-mage-init='{"tabs":{"openedState":"active"}}'>
                                            <div class="data item title"
                     aria-labeledby="tab-label-product.info.description-title"
                     data-role="collapsible" id="tab-label-product.info.description">
                    <a class="data switch"
                       tabindex="-1"
                       data-toggle="switch"
                       href="#product.info.description"
                       id="tab-label-product.info.description-title">
                        Overview                    </a>
                </div>
                <div class="data item content" id="product.info.description" data-role="content">
                    
<div class="overview-data">
    <div class="product attribute description">
                <div class="value layout-2-columns" >
                            <div class="column-left">
                    <!--Show Video or Image for Simple Product in here-->
                                        <!--Show Chart for Config Product in here-->
                    <div class="charts-list">
                                                    <div class="forix-chartjs"
                                 data-chart-background= "https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Bg-Cupping-02.png"
                                 data-chart-label=  "Papua New Guinea Honey Nebilyer Estate"
                                 data-chart-id=     "cupping-chart"
                                 data-chart-type=   "radar"
                                 data-chart-value=  "Dry Fragrance:8.5,Wet Aroma:8.7,Brightness:8,Flavor:8.7,Body:8.6,Finish:7.7,Sweetness:8.6,Clean Cup:7.3,Complexity:8.8,Uniformity:8.2"
                                 data-chart-score=  "86.6"
                                 data-cupper-correction="3.5"
                            >
                            </div>
                            <canvas id="cupping-chart" class="chart-area type-radar loading"></canvas>
                        
                                                    <div class="forix-chartjs"
                                 data-chart-background= "https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Bg-FlavorChart-01.png"
                                 data-chart-label=  "Papua New Guinea Honey Nebilyer Estate"
                                 data-chart-id=     "flavor-chart"
                                 data-chart-type=   "polarArea"
                                 data-chart-value=  "Floral:0,Honey:0,Sugars:3,Caramel:0,Fruits:3,Citrus:1,Berry:3,Cocoa:3,Nuts:0,Rustic:3,Spice:1,Body:3"
                            >
                            </div>
                            <canvas id="flavor-chart" class="chart-area type-polar loading"></canvas>
                                            </div>    
                </div>
                                        <div class="column-right">
                                            <div class="box-description">
                            <div class="list-info">
                                <ul class="list-info">
                                                                                                                        <li>
                                                <strong>Process Method</strong>
                                                <span>Honey Process</span>
                                            </li>
                                                                                                                                                                <li>
                                                <strong>Cultivar</strong>
                                                <span>Bourbon Types, Typica Types</span>
                                            </li>
                                                                                                                                                                <li>
                                                <strong>Farm Gate</strong>
                                                <span>No</span>
                                            </li>
                                                                                                                                                                                        </ul>
                            </div>
                                                    </div>
                                    </div>
                    </div>
    </div>
</div>
                </div>
                                            <div class="data item title"
                     aria-labeledby="tab-label-product.info.specs-title"
                     data-role="collapsible" id="tab-label-product.info.specs">
                    <a class="data switch"
                       tabindex="-1"
                       data-toggle="switch"
                       href="#product.info.specs"
                       id="tab-label-product.info.specs-title">
                        Specs                    </a>
                </div>
                <div class="data item content" id="product.info.specs" data-role="content">
                        <div class="additional-attributes-wrapper table-wrapper table-scroll">
        <table class="additional-attributes-table" id="product-attribute-specs-table">
            <tbody>
            <tr>
                        <th class="col label" scope="row">Region</th>
                        <td class="col data" data-th="Region">
                            Waghi Valley</tr><tr>
                        <th class="col label" scope="row">Processing</th>
                        <td class="col data" data-th="Processing">
                            Honey Process </tr><tr>
                        <th class="col label" scope="row">Drying Method</th>
                        <td class="col data" data-th="Drying Method">
                            Patio Sun-dried</tr><tr>
                        <th class="col label" scope="row">Arrival date</th>
                        <td class="col data" data-th="Arrival date">
                            November 2018 Arrival</tr><tr>
                        <th class="col label" scope="row">Lot size</th>
                        <td class="col data" data-th="Lot size">
                            32</tr><tr>
                        <th class="col label" scope="row">Bag size</th>
                        <td class="col data" data-th="Bag size">
                            60 KG</tr><tr>
                        <th class="col label" scope="row">Packaging</th>
                        <td class="col data" data-th="Packaging">
                            GrainPro liner</tr><tr>
                        <th class="col label" scope="row">Cultivar Detail</th>
                        <td class="col data" data-th="Cultivar Detail">
                            Arusha, Bourbon, Typica</tr><tr>
                        <th class="col label" scope="row">Grade</th>
                        <td class="col data" data-th="Grade">
                            A/X</tr><tr>
                        <th class="col label" scope="row">Appearance</th>
                        <td class="col data" data-th="Appearance">
                            .8 d per 300 grams, 15 - 19 Screen</tr><tr>
                        <th class="col label" scope="row">Roast Recommendations</th>
                        <td class="col data" data-th="Roast Recommendations">
                            City+ to Full City+</tr><tr>
                        <th class="col label" scope="row">Recommended for Espresso</th>
                        <td class="col data" data-th="Recommended for Espresso">
                            Yes</tr>            </tbody>
        </table>
    </div>
                </div>
                                                            <div class="data item title"
                     aria-labeledby="tab-label-product-info-origin-notes-title"
                     data-role="collapsible" id="tab-label-product-info-origin-notes">
                    <a class="data switch"
                       tabindex="-1"
                       data-toggle="switch"
                       href="#product-info-origin-notes"
                       id="tab-label-product-info-origin-notes-title">
                        Farm Notes                    </a>
                </div>
                <div class="data item content" id="product-info-origin-notes" data-role="content">
                    
    <div class="product attribute origin-notes">
        <div class="value layout-2-columns">
            <div class="column-left">
                                                        <div id="circle1" class="map-circle"
                         data-lat= "-5.870922" data-lng= "144.644204"
                         data-radius="2000" data-color="15z">
                    </div>
                                                </div>
                            <div class="column-right">
                    <p>This coffee comes to us by way of the Kuta coffee mill in the Waghi District of Papua New Guinea. The coffee processed at the mill are from smaller coffee plantations in the area situated at just under 1600 meters above sea level on the low end. It's a honey processed coffee, meaning the coffee cherry and much of the fruit are stripped from the seed using depulping machinery, and then the seed still covered in sticky mucilage is laid to dry with any remaining fruit still intact. This tends to result in bigger body, softer acidity, and often a fruited cup. The physical grade of this coffee is impressive, and I couldn't find a single full quaker bean in the few hundred grams of coffee I roasted.</p>                </div>
                    </div>
    </div>
    <script async defer src="https://maps.googleapis.com/maps/api/js?key=AIzaSyAIM3wNOD-JZKwba-4L0f_E0_rEiKRcAsI" type="text/javascript"></script>



                </div>
                                            <div class="data item title"
                     aria-labeledby="tab-label-product-info-cupping-notes-title"
                     data-role="collapsible" id="tab-label-product-info-cupping-notes">
                    <a class="data switch"
                       tabindex="-1"
                       data-toggle="switch"
                       href="#product-info-cupping-notes"
                       id="tab-label-product-info-cupping-notes-title">
                        Cupping Notes                    </a>
                </div>
                <div class="data item content" id="product-info-cupping-notes" data-role="content">
                    <div class="product attribute cupping-notes">
    <div class="value layout-2-columns">

                    <div class="chart-cupping-notes column-left">
                <div class="forix-chartjs"
                     data-chart-background= "https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Bg-Cupping-02.png"
                     data-chart-label=  "Papua New Guinea Honey Nebilyer Estate"
                     data-chart-id=     "cupping-chart2"
                     data-chart-type=   "radar"
                     data-chart-value=  "Dry Fragrance:8.5,Wet Aroma:8.7,Brightness:8,Flavor:8.7,Body:8.6,Finish:7.7,Sweetness:8.6,Clean Cup:7.3,Complexity:8.8,Uniformity:8.2"
                     data-chart-score=  "86.6"
                     data-cupper-correction="3.5"
                >
                </div>
                <canvas id="cupping-chart2" class="chart-area type-radar loading"></canvas>
            </div>
        
                    <div class="column-right">
                                    <p>A honey process batch from the same mill that brought us Kuta Waghi (we still have some available). In fact, tasting these two coffees side by side, you get an idea of the role processing plays in a coffee's final cup profile. This lot is much more fruited than the Peaberry, and the aroma displays a sweet blend of honey and molasses sweetness, and fruited smells like dark berry pulp drifting up in the steam. The cup displays a nice fruited profile as well when roasted to City+, and a sort of unrefined sweetness at the core helps to highlight fruited nuance. The cooling coffee reveals glimpses of red grape, raisin, and a dark cranberry note in the finish. Full City roasts are also highlighted by berry characteristics, along with a rustic dark chocolate base note. Acidity is quite mild (typical of honey process coffee) across the roast spectrum, and Full City roasts produce inky espresso shots with layers deep cocoa roast tone interspersed with a sweet cranberry hint.</p>                                            </div>
        

    </div>
</div>
                </div>
                    </div>
    </div>
<script type="text/x-magento-init">
    {
        "*": {
            "forix/chart": {}
        }
    }
</script><input name="form_key" type="hidden" value="2hKQSKcftqnlC4Bx" /><div id="authenticationPopup" data-bind="scope:'authenticationPopup'" style="display: none;">
    <script>
        window.authenticationPopup = {"customerRegisterUrl":"https:\/\/www.sweetmarias.com\/customer\/account\/create\/","customerForgotPasswordUrl":"https:\/\/www.sweetmarias.com\/customer\/account\/forgotpassword\/","baseUrl":"https:\/\/www.sweetmarias.com\/"};
    </script>
    <!-- ko template: getTemplate() --><!-- /ko -->
    <script type="text/x-magento-init">
        {
            "#authenticationPopup": {
                "Magento_Ui/js/core/app": {"components":{"authenticationPopup":{"component":"Magento_Customer\/js\/view\/authentication-popup","children":{"messages":{"component":"Magento_Ui\/js\/view\/messages","displayArea":"messages"},"captcha":{"component":"Magento_Captcha\/js\/view\/checkout\/loginCaptcha","displayArea":"additional-login-form-fields","formId":"user_login","configSource":"checkout"}}}}}            },
            "*": {
                "Magento_Ui/js/block-loader": "https://sweet-marias-live-sweetmarias.netdna-ssl.com/static/version1553065162/frontend/Forix/sweetmarias/en_US/images/loader-1.gif"
            }
        }
    </script>
</div>
<script type="text/x-magento-init">
{"*":{"Magento_Customer\/js\/section-config":{"sections":{"stores\/store\/switch":"*","directory\/currency\/switch":"*","*":["messages"],"customer\/account\/logout":"*","customer\/account\/loginpost":"*","customer\/account\/createpost":"*","customer\/ajax\/login":["checkout-data","cart","customer","compare-products","wishlist"],"catalog\/product_compare\/add":["compare-products"],"catalog\/product_compare\/remove":["compare-products"],"catalog\/product_compare\/clear":["compare-products"],"giftcard\/cart\/add":["cart"],"giftcard\/cart\/remove":["cart"],"mlogin\/account\/loginpost":["customer","checkout-data","cart","compare-products","wishlist"],"sales\/guest\/reorder":["cart"],"sales\/order\/reorder":["cart"],"checkout\/cart\/add":["cart"],"checkout\/cart\/delete":["cart"],"checkout\/cart\/updatepost":["cart"],"checkout\/cart\/updateitemoptions":["cart"],"checkout\/cart\/couponpost":["cart"],"checkout\/cart\/estimatepost":["cart"],"checkout\/cart\/estimateupdatepost":["cart"],"checkout\/onepage\/saveorder":["cart","checkout-data","last-ordered-items"],"checkout\/sidebar\/removeitem":["cart"],"checkout\/sidebar\/updateitemqty":["cart"],"rest\/*\/v1\/carts\/*\/payment-information":["cart","checkout-data","last-ordered-items"],"rest\/*\/v1\/guest-carts\/*\/payment-information":["cart","checkout-data"],"rest\/*\/v1\/guest-carts\/*\/selected-payment-method":["cart","checkout-data"],"rest\/*\/v1\/carts\/*\/selected-payment-method":["cart","checkout-data"],"review\/product\/post":["review"],"wishlist\/index\/add":["wishlist"],"wishlist\/index\/remove":["wishlist"],"wishlist\/index\/updateitemoptions":["wishlist"],"wishlist\/index\/update":["wishlist"],"wishlist\/index\/cart":["wishlist","cart"],"wishlist\/index\/fromcart":["wishlist","cart"],"wishlist\/index\/allcart":["wishlist","cart"],"wishlist\/shared\/allcart":["wishlist","cart"],"wishlist\/shared\/cart":["cart"],"giftregistry\/index\/cart":["cart"],"giftregistry\/view\/addtocart":["cart"],"wishlist\/index\/copyitem":["wishlist"],"wishlist\/index\/copyitems":["wishlist"],"wishlist\/index\/deletewishlist":["wishlist","multiplewishlist"],"wishlist\/index\/createwishlist":["multiplewishlist"],"wishlist\/index\/moveitem":["wishlist"],"wishlist\/index\/moveitems":["wishlist"],"wishlist\/search\/addtocart":["cart","wishlist"],"multishipping\/checkout\/overviewpost":["cart"],"authorizenet\/directpost_payment\/place":["cart","checkout-data"],"customer_order\/cart\/updatefaileditemoptions":["cart"],"checkout\/cart\/updatefaileditemoptions":["cart"],"customer_order\/cart\/advancedadd":["cart"],"checkout\/cart\/advancedadd":["cart"],"checkout\/cart\/removeallfailed":["cart"],"customer_order\/cart\/addfaileditems":["cart"],"checkout\/cart\/addfaileditems":["cart"],"customer_order\/sku\/uploadfile":["cart"],"paypal\/express\/placeorder":["cart","checkout-data"],"paypal\/payflowexpress\/placeorder":["cart","checkout-data"],"braintree\/paypal\/placeorder":["cart","checkout-data"]},"clientSideSections":["checkout-data"],"baseUrls":["https:\/\/www.sweetmarias.com\/"]}}}</script>
<script type="text/x-magento-init">
{"*":{"Magento_Customer\/js\/customer-data":{"sectionLoadUrl":"https:\/\/www.sweetmarias.com\/customer\/section\/load\/","cookieLifeTime":"180000","updateSessionUrl":"https:\/\/www.sweetmarias.com\/customer\/account\/updateSession\/"}}}</script>
<script type="text/x-magento-init">
    {
        "body": {
            "requireCookie": {"noCookieUrl":"https:\/\/www.sweetmarias.com\/cookie\/index\/noCookies\/","triggers":{"addToWishlistLink":".action.towishlist"}}        }
    }
</script>
<script type="text/x-magento-init">
    {
        "body": {
            "pageCache": {"url":"https:\/\/www.sweetmarias.com\/page_cache\/block\/render\/id\/14051\/","handles":["default","catalog_product_view","catalog_product_view_id_14051","catalog_product_view_sku_GCX-6048","catalog_product_view_type_configurable"],"originalRequest":{"route":"catalog","controller":"product","action":"view","uri":"\/papua-new-guinea-honey-process-nebilyer-estate.html"},"versionCookieName":"private_content_version"}        }
    }
</script>
<div id="monkey_campaign" style="display:none;"
     data-mage-init='{"campaigncatcher":{"checkCampaignUrl": "https://www.sweetmarias.com/mailchimp/campaign/check/"}}'>
</div>

</div></div></main><footer class="page-footer"><div class="footer top"><div class="footer-top inner"><ul class="footer_links">
<li class="footer-about-us">
<h3 class="title">About Us</h3>
<ul class="content">
<li>
<p>(510)628-0919 <br />(888)876-5917 Toll Free <br />info@sweetmarias.com<br />Mon-Fri 10:00am-5:00pm PST</p>
</li>
<li class="link-bottom-menu"><a href="https://www.sweetmarias.com/about-us/" target="_self">About Us</a></li>
<li class="link-bottom-menu"><a href="https://www.sweetmarias.com/contact/" target="_self">Contact Us</a></li>
</ul>
</li>
<li class="footer-help">
<h3 class="title">Help</h3>
<ul class="content"><!--<li><a href="https://www.sweetmarias.com/holidays" target="_blank">Holidays</a>&nbsp;</li><-->
<li><a title="Ordering" href="https://www.sweetmarias.com/ordering-policies/" target="_self">Ordering</a></li>
<li><a title="Shipping" href="https://www.sweetmarias.com/shipping-policies/" target="_self">Shipping</a></li>
<li><a title="Returns" href="https://www.sweetmarias.com/returns-policies/" target="_self">Returns</a></li>
<li class="link-bottom-menu"><a title="FAQ" href="https://www.sweetmarias.com/faq/" target="_self">FAQs</a></li>
</ul>
</li>
<li class="footer-roasting-basic">
<h3 class="title">Roasting Basics</h3>
<ul class="content">
<li><a href="https://legacy.sweetmarias.com/library/category/roast/getting-started/" target="_self">Home Roasting Basics</a></li>
<li><a href="https://legacy.sweetmarias.com/library/choosing-a-roaster-part-1/" target="_self">Choosing a Roaster</a></li>
<li><a href="https://legacy.sweetmarias.com/library/choosing-green-coffee-replacing-a-favorite/" target="_self">Choosing Green Coffee</a></li>
<li><a href="https://legacy.sweetmarias.com/library/storing-your-roasted-coffee-2/" target="_self">Coffee Storage</a></li>
<li><a href="https://legacy.sweetmarias.com/library/using-sight-to-determine-degree-of-roast/" target="_self">Roast Process Pictorial</a></li>
</ul>
</li>
<li class="footer-resources">
<h3 class="title">Resources</h3>
<ul class="content">
<li><a title="Our Coffee Library" href="https://legacy.sweetmarias.com/library/" target="_self">Our Coffee Library</a></li>
<li><a title="Product Video" href="https://legacy.sweetmarias.com/library/category/video/" target="_self">Product Video</a></li>
<li><a title="Coffee Glossary" href="https://www.sweetmarias.com/glossary/" target="_self">Coffee Glossary</a></li>
</ul>
</li>
</ul><div class="footer_visit"><a href="https://coffeeshrub.com/" target="_blank"> <img src="https://sweet-marias-live-sweetmarias.netdna-ssl.com/media/wysiwyg/Homepage/coffee-shrub-wholesale-logo.jpg" alt="Visit Coffee Shrub" width="139" height="128" /> </a>
<h5>Coffee Shrub: Sweet Maria's Wholesale Site</h5>
</div></div></div><div class="footer content"><div class="footer inner"><div class="block newsletter">
    <h3 class="title">Never Miss An Update</h3>
    <div class="content">
        <form class="form subscribe"
            novalidate
            action="https://www.sweetmarias.com/newsletter/subscriber/new/"
            method="post"
            data-mage-init='{"validation": {"errorClass": "mage-error"}}'
            id="newsletter-validate-detail">
            <div class="field newsletter">
            <!--<label class="label" for="newsletter"><span>--><!--</span></label>-->
                <div class="control">
                    <input name="email" type="email" id="newsletter"
                                placeholder="Your email address"
                                data-validate="{required:true, 'validate-email':true}"/>
                </div>
            </div>
            <div class="actions">
                <button class="action subscribe" title="Subscribe" type="submit">
                    <span>Subscribe</span>
                </button>
            </div>
        </form>
    </div>
</div>

<div class="block socials">
    <h3 class="title">Stay In Touch</h3>
    <ul class="box-socials">
        <li class="icon-facebook">
            <a href="https://www.facebook.com/SweetMariasCoffee" target="_blank" rel="nofollow" title="Facebook"><span>Facebook</span></a>
        </li>
        <li class="icon-twitter">
            <a href="https://twitter.com/sweetmarias" target="_blank" rel="nofollow" title="Twitter"><span>Twitter</span></a>
        </li>
        <li class="icon-youtube">
            <a href="https://www.youtube.com/user/sweetmarias" target="_blank" rel="nofollow" title="Youtube"><span>Youtube</span></a>
        </li>
        <li class="icon-instagram">
            <a href="https://www.instagram.com/sweetmarias/" target="_blank" rel="nofollow" title="Instagram"><span>Instagram</span></a>
        </li>
    </ul>
</div>
</div></div><div class="footer bottom"><div class="footer-bottom inner"><small class="copyright">
    <span>©2019 Sweet Maria’s. All rights reserved.</span>
</small>
</div></div></footer>
<script type="text/x-magento-init">
    {
        "*": {
            "forix/pdp": {}
        }
    }
</script>




<!-- Go to www.addthis.com/dashboard to customize your tools -->
<script type="text/javascript" src="//s7.addthis.com/js/300/addthis_widget.js#pubid=ra-5af479378586f56a"></script>
</div>    </body>
</html>
		`)
	}))
	defer ts.Close()

	c, err := LoadCoffee(ts.URL)
	if err != nil {
		t.Fatalf("Failed to LoadCoffee: %s", err)
	}

	assert.Equal(t, Coffee{
		Title:    "Papua New Guinea Honey Nebilyer Estate",
		Overview: "Honey process seems to bring out fruited notes like cranberry, raisin, red grape, and underscored by molasses-like sweetness. This PNG boasts body, and with mild acidity, is great espresso too. City+ to Full City+. Good for espresso.",
		Score:    86.6,
		URL:      ts.URL,

		FarmNotes:    "This coffee comes to us by way of the Kuta coffee mill in the Waghi District of Papua New Guinea. The coffee processed at the mill are from smaller coffee plantations in the area situated at just under 1600 meters above sea level on the low end. It's a honey processed coffee, meaning the coffee cherry and much of the fruit are stripped from the seed using depulping machinery, and then the seed still covered in sticky mucilage is laid to dry with any remaining fruit still intact. This tends to result in bigger body, softer acidity, and often a fruited cup. The physical grade of this coffee is impressive, and I couldn't find a single full quaker bean in the few hundred grams of coffee I roasted.",
		CuppingNotes: "A honey process batch from the same mill that brought us Kuta Waghi (we still have some available). In fact, tasting these two coffees side by side, you get an idea of the role processing plays in a coffee's final cup profile. This lot is much more fruited than the Peaberry, and the aroma displays a sweet blend of honey and molasses sweetness, and fruited smells like dark berry pulp drifting up in the steam. The cup displays a nice fruited profile as well when roasted to City+, and a sort of unrefined sweetness at the core helps to highlight fruited nuance. The cooling coffee reveals glimpses of red grape, raisin, and a dark cranberry note in the finish. Full City roasts are also highlighted by berry characteristics, along with a rustic dark chocolate base note. Acidity is quite mild (typical of honey process coffee) across the roast spectrum, and Full City roasts produce inky espresso shots with layers deep cocoa roast tone interspersed with a sweet cranberry hint.",

		AdditionalAttributes: map[string]string{
			"Appearance":               ".8 d per 300 grams, 15 - 19 Screen",
			"Arrival date":             "November 2018 Arrival",
			"Bag size":                 "60 KG",
			"Cultivar Detail":          "Arusha, Bourbon, Typica",
			"Drying Method":            "Patio Sun-dried",
			"Grade":                    "A/X",
			"Lot size":                 "32",
			"Packaging":                "GrainPro liner",
			"Processing":               "Honey Process",
			"Recommended for Espresso": "Yes",
			"Region":                   "Waghi Valley",
			"Roast Recommendations":    "City+ to Full City+",
		},
	}, c)
}
