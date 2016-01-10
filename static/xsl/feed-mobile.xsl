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
        <div class="container"> 
          <h1>Brooklet</h1>
          <xsl:apply-templates select="navigation" />
        </div>

        <div class="container"> 
          <xsl:apply-templates select="content" />
        </div>

      </body>

    </html>
  </xsl:template>

  <xsl:template match="navigation">
    <div class="mobile">
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
    </div>
  </xsl:template>

  <xsl:template match="entry">
    <div class="row">
      <div class="col mobile portrait">
        <a href="{url}">
        
        <xsl:choose>
          <xsl:when test="string-length(image/@src) > 0">
            <img src="/images/{image/@src}-512.png">Test</img>
          </xsl:when>
          <xsl:otherwise>
            <xsl:variable name="varcolor"><xsl:value-of select="string-length(title) mod 9 + 1" /></xsl:variable>
            <div class="colordiv color{$varcolor}"></div>
          </xsl:otherwise>
        </xsl:choose>

        <div class="textbox">
          <h1><xsl:apply-templates select="title" /></h1>
        </div>
        </a>
      </div>
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
