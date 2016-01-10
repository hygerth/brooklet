<?xml version="1.0"?>
<xsl:stylesheet version="1.0"
    xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns="http://www.w3.org/1999/xhtml">

    <xsl:template match="page">
        <html>
            <head>
                <title>Brooklet</title>
                <link href="/static/css/screen.css" rel="stylesheet" type="text/css" />
            </head>

            <body>
                <header> 
                    <h1>Brooklet</h1>
                    <xsl:apply-templates select="navigation" />
                </header>

                <div class="container"> 
                    <ol>
                        <li>
                            <h2>Add new feed</h2>
                            <form action="/add/feed" method="POST">
                                <input type="text" name="feed" placeholder="url"/>
                                <input type="submit" value="Add" />
                            </form>
                        </li>
                       <h2>Current Feeds </h2>
                        <xsl:for-each select="subscription">
                            <li>
                                    <p><form action="/remove/feed" method="POST">
                                        <input type="hidden" name="feed" value="{url}" />
                                        <input type="submit" value="&#x2716;"/>
                                    </form>
                                    <xsl:value-of select="title" /></p>
                            </li>
                        </xsl:for-each>
                    </ol>

                    <h2>Filter Feed Keywords</h2>
                    <ol>
                        <xsl:for-each select="filter">
                            <li>
                                <p>
                                <form action="/remove/filter" method="POST">
                                    <input type="hidden" name="filter" value="{.}" />
                                    <input type="submit" value="&#x2716;"/>
                                </form>
                                <xsl:value-of select="." />
                                </p>
                            </li>
                        </xsl:for-each>
                        <li>
                            <form action="/add/filter" method="POST">
                                <input type="text" name="filter" placeholder="keyword"/>
                                <input type="submit" value="Add" />
                            </form>
                        </li>
                    </ol>
                </div>

            </body>

        </html>
    </xsl:template>

    <xsl:template match="navigation">
        <label for="show-menu" class="show-menu">Show Menu</label>
        <input type="checkbox" id="show-menu" role="button"/>
        <ul class="menu">
            <xsl:for-each select="navigationitem">
                <li>
                    <a href="{@url}">
                        <xsl:value-of select="@title" />
                    </a>
                    <ul class="hidden">
                        <xsl:for-each select="sublist/item">
                            <li>
                                <a href="{url}">
                                    <xsl:apply-templates select="title" />
                                </a>
                            </li>
                        </xsl:for-each>
                    </ul>
                </li>
            </xsl:for-each>
        </ul>
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
