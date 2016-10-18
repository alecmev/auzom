import { isEqual } from 'lodash';
import { PropTypes } from 'react';
import ReactMarkdown from 'react-markdown';
import sanitizeHtml from 'sanitize-html';
import slugify from 'slug';

import * as utils from '../utils';
import Component from '../utils/Component';

// markdown norm:
// 'blockquote', 'em', 'br', 'img', 'a', 'p', 'strong', 'hr',
// 'ul', 'ol', 'li', 'code', 'pre',
// 'h1', 'h2', 'h3', 'h4', 'h5', 'h6',

const sanitizationPolicy = {
  allowedTags: [
    'em', 'br', 'a', 'img', 'strong',
    'table', 'thead', 'tbody', 'tfoot', 'tr', 'th', 'td',
    'caption', 'colgroup', 'col', // exotic, but whatever
  ],
  allowedAttributes: {
    a: ['href', 'name', 'target'],
    img: ['src'],
    colgroup: ['span'],
    col: ['span'],
    th: ['colspan', 'rowspan'],
    td: ['colspan', 'rowspan'],
  },
};

function htmlRenderer(x) {
  return React.createElement(x.isBlock ? 'div' : 'span', {
    dangerouslySetInnerHTML: {
      __html: sanitizeHtml(x.literal, sanitizationPolicy),
    },
  });
}

const renderers = {
  HtmlBlock: htmlRenderer,
  HtmlInline: htmlRenderer,
  Heading: x => React.createElement(`h${x.level}`, null, [
    <div key={x.nodeKey} className="u-anchor" id={x.literal} />,
    x.children,
  ]),
};

function ast2text(node) {
  const walker = node.walker();
  let event;
  let acc = '';
  while (event = walker.next()) { // eslint-disable-line
    if (event.node.literal) acc += event.node.literal;
  }

  return acc;
}

export default class Markdown extends Component {
  static propTypes = {
    source: PropTypes.string,
    onToc: PropTypes.func,
  };

  walker = (x) => {
    if (x.node.type === 'Heading') {
      if (!x.entering) return;
      const text = ast2text(x.node);
      const slug = slugify(text, { mode: 'rfc3986' });
      let id = slug;
      if (this.tocSlugs[slug]) id += `-${this.tocSlugs[slug]}`;
      this.tocSlugs[slug] = this.tocSlugs[slug] + 1 || 1;
      // this is a hack; only literal is available in the heading renderer
      x.node.literal = id;
      const item = { id, level: x.node.level, text };
      this.toc.push(item);
      this.tocDiff = this.tocDiff || !(
        this.tocOld &&
        this.tocOld.length >= this.toc.length &&
        isEqual(this.tocOld[this.toc.length - 1], item)
      );
    } else if (x.node.type === 'Document') {
      if (x.entering) {
        this.tocOld = this.toc;
        this.toc = [];
        this.tocDiff = !this.tocOld;
        this.tocSlugs = {};
      } else if (this.tocDiff && this.props.onToc) this.props.onToc(this.toc);
    }
  };

  render() {
    return (
      <ReactMarkdown
        className={this.cni()}
        renderers={renderers}
        transformImageUri={utils.https}
        walker={this.walker}
        source={this.props.source}
      />
    );
  }
}
