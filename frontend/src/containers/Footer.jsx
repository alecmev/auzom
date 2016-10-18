import { Link } from 'react-router';

import Component from '../utils/Component';

export default class Footer extends Component {
  render() {
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cni({ u: 'sectionMargined' })}>
          <img
            className={this.cn({ d: 'logo' })}
            src="//i.imgur.com/xUXpADb.png"
          />
          <div className={this.cn({ d: 'who' })}>
            <div className={this.cn({ d: 'madeBy' })}>made by</div>
            <div className={this.cn({ d: 'auzom' })}>auzom</div>
            <div className={this.cn({ d: 'entertainment' })}>entertainment</div>
          </div>
          <div className={this.cn({ d: 'filler' })} />
          <Link
            className={this.cn({ d: 'link' })}
            to="/about"
          >
            About
          </Link>
          <a
            className={this.cn({ d: 'link' })}
            href="//blog.auzom.gg"
          >
            Blog
          </a>
          <a
            className={this.cn({ d: 'link' })}
            href="mailto:support@auzom.gg"
          >
            Contact us
          </a>
          <a
            className={this.cn({ d: 'social', m: 'first' })}
            href="//twitch.tv/team/auzom"
          >
            <img src={require('../assets/twitch.svg')} />
          </a>
          <a
            className={this.cn({ d: 'social' })}
            href="//youtube.com/c/auzom_gg"
          >
            <img src={require('../assets/youtube.svg')} />
          </a>
          <a
            className={this.cn({ d: 'social' })}
            href="//twitter.com/auzom_gg"
          >
            <img src={require('../assets/twitter.svg')} />
          </a>
          <a
            className={this.cn({ d: 'social' })}
            href="//facebook.com/auzom.gg"
          >
            <img src={require('../assets/facebook.svg')} />
          </a>
        </div>
      </div>
    );
  }
}
