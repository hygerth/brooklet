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
                    <h1>Subscriptions</h1>
                    <ul>
                        <xsl:for-each select="subscription">
                            <li>
                                <a href="/feed/{id}">
                                    <xsl:value-of select="title" />
                                </a>
                            </li>
                        </xsl:for-each>
                    </ul>
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

    <xsl:template match="title">
        <xsl:apply-templates />
    </xsl:template>

    <xsl:template match="url">
        <xsl:apply-templates />
    </xsl:template>

    <xsl:template match="published">
        <xsl:apply-templates />
    </xsl:template>

    <xsl:template match="summary">
        <xsl:apply-templates />
    </xsl:template>

    <xsl:template match="author">
        <xsl:apply-templates />
    </xsl:template>

    <xsl:template match="name">
        <xsl:apply-templates />
    </xsl:template>

    <xsl:template match="url">
        <xsl:apply-templates />
    </xsl:template>

</xsl:stylesheet>
