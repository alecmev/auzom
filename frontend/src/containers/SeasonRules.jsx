import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { createStructuredSelector } from 'reselect';

import * as selectors from '../selectors';
import Component from '../utils/Component';

import Markdown from '../components/Markdown';

@connect(createStructuredSelector({
  amAdmin: selectors.amAdmin,
}))
export default class SeasonRules extends Component {
  static propTypes = {
    location: PropTypes.object, // eslint-disable-line
    season: ImmutablePropTypes.map,
    amAdmin: PropTypes.bool,
  };

  state = {
    toc: null,
    hashChecked: false,
  };

  onToc = (toc) => {
    // To avoid setState being called during render and to also let the element
    // we need to scroll to be added to DOM first.
    setTimeout(() => {
      this.setState({ toc: toc.map(x => (
        <a
          key={x.id}
          className={this.cn({ d: 'tocItem', m: x.level })}
          href={`#${x.id}`}
        >
          {x.text}
        </a>
      )) });
      if (this.state.hashChecked) return;
      this.setState({ hashChecked: true });
      if (location.hash.length) {
        const x = document.getElementById(location.hash.split('#')[1]);
        x && x.scrollIntoView();
      }
    }, 0);
  }

  render() {
    const { season, amAdmin } = this.props;
    if (!season) return null;
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ u: 'sectionMargined', d: 'inner' })}>
          <div className={this.cn({ d: 'content' })}>
            <Markdown source={season.get('rules')} onToc={this.onToc} />
          </div>
          <div className={this.cn({ d: 'toc' })}>
            {amAdmin && <Link
              className={this.cn({ d: 'tocItem', m: ['1', 'edit'] })}
              to={`${this.props.location.pathname}/edit`}
            >EDIT ✎</Link>}
            {this.state.toc}
          </div>
        </div>
      </div>
    );
  }
}
