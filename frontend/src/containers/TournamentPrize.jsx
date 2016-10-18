import { PropTypes } from 'react';

import Component from '../utils/Component';

export default class TournamentPrize extends Component {
  static propTypes = {
    item: PropTypes.object.isRequired,
  };

  render() {
    const { item } = this.props;
    return (
      <a
        className={this.cni({ m: item.isSmall && 'small' })}
        href={item.url}
      >
        <img className={this.cn({ d: 'picture' })} src={item.picture} />
        <div className={this.cn({ d: 'glow', m: item.metal })} />
        <div className={this.cn({ d: 'header' })}>
          <span className={this.cn({ d: 'place', m: item.metal })}>
            {item.place}
          </span>
          <span className={this.cn({ d: 'scope' })}>
            in {item.scope}
          </span>
          <span className={this.cn({ d: 'pieces' })}>
            {item.pieces}x
          </span>
        </div>
        <div className={this.cn({ d: 'name' })}>{item.name}</div>
      </a>
    );
  }
}
