<?xml version="1.0"?>
<xsl:stylesheet version="1.0"
    xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns="http://www.w3.org/1999/xhtml">

    <xsl:template match="page">
        <html>
            <head>
                <title>Brooklet</title>
                <link href="/static/css/screen.css" rel="stylesheet" type="text/css" />
                <meta name="viewport" content="initial-scale=1.0"/>
            </head>

            <body>
                <header> 
                    <xsl:apply-templates select="navigation" />
                </header>

                <div class="container"> 
                    <xsl:apply-templates select="content" />
                </div>

            </body>

        </html>
    </xsl:template>

    <xsl:template match="navigation">
        <div class="mobile">
            <input type="checkbox" id="nav-trigger" class="nav-trigger" />
            <label for="nav-trigger"><span></span></label>
            <nav>
                <h1>B</h1>
                <ul>
                    <xsl:for-each select="navigationitem">
                        <li>
                            <a href="{@url}">
                                <xsl:value-of select="@title" />
                            </a>
                            <ul>
                                <xsl:for-each select="sublist/item">
                                    <li>
                                        <a href="{url}">
                                            <xsl:value-of select="title" />
                                        </a>
                                    </li>
                                </xsl:for-each>
                            </ul>
                        </li>
                    </xsl:for-each>
                </ul>
            </nav>
        </div>
    </xsl:template>

    <xsl:template match="content">
        <xsl:apply-templates select="entry" />
    </xsl:template>

    <xsl:template match="entry">
        <div class="row">
            <div class="col mobile portrait">
                <a href="/article/{id}">
        
                    <xsl:choose>
                        <xsl:when test="not(not(hasimage='true'))">
                            <img src="/images/{id}-128.png"/>
                        </xsl:when>
                        <xsl:otherwise>
                            <xsl:variable name="varcolor"><xsl:value-of select="string-length(title) mod 9 + 1" /></xsl:variable>
                            <div class="colordiv color{$varcolor}"></div>
                        </xsl:otherwise>
                    </xsl:choose>

                    <div class="textbox">
                        <h1><xsl:value-of select="title" /></h1>
                    </div>
                </a>
            </div>
        </div>
    </xsl:template>
</xsl:stylesheet>
