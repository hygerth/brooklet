<?xml version="1.0"?>
<xsl:stylesheet version="1.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
  xmlns="http://www.w3.org/1999/xhtml">

  <xsl:template match="page">
    <html>
      <head>
        <title>Brooklet</title>
        <link href="style.css" rel="stylesheet" type="text/css" />
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

  <xsl:template match="entry[position() mod 2 = 1]">
    <div class="row">
      <xsl:apply-templates mode="proc" select=".|following-sibling::entry[not(position() > 1)]" />
    </div>
  </xsl:template>

  <xsl:template match="entry" mode="proc">
      <div class="col {image/@rotation}">
        <a href="{url}">
        <img src="{image}">Test</img>
        <div class="textbox">
          <h1><xsl:apply-templates select="title" /></h1>
        </div>
        </a>
      </div>
  </xsl:template>

  <xsl:template match="entry[not(position() mod 2 = 1)]" />

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
