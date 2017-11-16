import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { createStructuredSelector } from 'reselect';

import * as selectors from '../selectors';
import Component from '../utils/Component';

@connect(createStructuredSelector({ me: selectors.me }))
export default class Home extends Component {
  static propTypes = {
    me: ImmutablePropTypes.map,
  };

  render() {
    const { me } = this.props;
    return (
      <div className={this.cni()}>
        <div className={this.cn({ u: 'sectionSingle', d: 'hero' })}>
          <div className={this.cn({ u: 'sectionMargined', d: 'heroInner' })}>
            <div className={this.cn({ d: 'heroFiller' })} />
            <div className={this.cn({ d: 'slogan' })}>
              <div>
                <span className={this.cn({ d: 'sloganBold' })}>
                  community
                </span> <span className={this.cn({ d: 'sloganNormal' })}>
                  soul
                </span>
              </div>
              <div className={this.cn({ d: 'sloganBottom' })}>
                <span className={this.cn({ d: 'sloganNormal' })}>
                  + pro
                </span> <span className={this.cn({ d: 'sloganBold' })}>
                  quality
                </span>
              </div>
            </div>
            <div className={this.cn({ d: 'heroSplitter' })} />
            <div className={this.cn({ d: 'youare' })}>
              Are you an aspiring competitive player?<br />
              Or are you a normal gamer looking for a challenge?<br />
              Or do you maybe want to take your cup to the next level?<br />
              <div className={this.cn({ d: 'youarePunch' })}>
                Whoever you are, we have something for you.
              </div>
            </div>
            <div className={this.cn({ d: 'heroFiller' })} />
          </div>
        </div>
        <div className={this.cni({ u: 'sectionSingle', d: 'row' })}>
          <div className={this.cn({ u: 'sectionMargined', d: 'rowInner' })}>
            <Link
              className={this.cn({ d: 'rowSidebar' })}
              to="/battlefield-4/conquest-league"
            >
              <div className={this.cn({ d: 'rowSidebarInner' })}>
                BCL
              </div>
            </Link>
            <div className={this.cn({ d: 'rowAbout' })}>
              <div className={this.cn({ d: 'rowTitle' })}>
                Battlefield Conquest League
              </div>
              <div className={this.cn({ d: 'rowInfo' })}>

                Being one of the oldest and largest Battlefield leagues out
                there, BCL dominates the European Battlefield 4 competitive
                scene. Even though Battlefield 4 is nearing its end, most recent
                season has attracted 54 teams, amounting to almost 1000 players

              </div>
              <Link
                className={this.cn({ d: 'rowLink' })}
                to="/battlefield-4/conquest-league"
              >
                Visit the tournament page →
              </Link>
            </div>
          </div>
        </div>
        <div className={this.cni({ u: 'sectionSingle', d: 'row', m: 'right' })}>
          <div className={this.cn({ u: 'sectionMargined', d: 'rowInner' })}>
            <div className={this.cn({ d: 'rowAbout' })}>
              <div className={this.cn({ d: 'rowTitle' })}>
                Aces High
              </div>
              <div className={this.cn({ d: 'rowInfo' })}>

                A community of players passionate about air combat in
                Battlefield 4. Aces High has been around for a couple of years
                now, focusing on training newcomers and mastering their assets,
                and after running a couple of test cups they’ve finally decided
                to up their game

              </div>
              <Link
                className={this.cn({ d: 'rowLink' })}
                to="/battlefield-4/aces-high-ah"
              >
                Visit the tournament page →
              </Link>
            </div>
            <Link
              className={this.cn({ d: 'rowSidebar' })}
              to="/battlefield-4/aces-high-ah"
            >
              <div className={this.cn({ d: 'rowSidebarInner' })}>
                ACE
              </div>
            </Link>
          </div>
        </div>
        <div className={this.cni({ u: 'sectionSingle', d: 'row', m: 'last' })}>
          <div className={this.cn({ u: 'sectionMargined', d: 'rowInner' })}>
            <div className={this.cn({ d: 'rowSidebar' })}>
              <div className={this.cn({ d: 'rowSidebarInner' })}>
                <div className={this.cn({ d: 'rowSidebarStat' })}>
                  1900+
                  <div>PLAYERS</div>
                </div>
                <div className={this.cn({ d: 'rowSidebarStat' })}>
                  200+
                  <div>TEAMS</div>
                </div>
                <div className={this.cn({ d: 'rowSidebarStat' })}>
                  490+
                  <div>MATCHES</div>
                </div>
              </div>
            </div>
            <div className={this.cn({ d: 'rowAbout' })}>
              <div className={this.cn({ d: 'rowTitle' })}>
                So... what is this all about?
              </div>
              <div className={this.cn({ d: 'rowInfo' })}>

                In a nutshell, we are a tournament platform with some novel
                extras.<br /><br />

                Our bread and butter are community-based competitions. They come
                to us and we nourish them into professional-grade organizations.
                Whatever the issues are - lacking exposure, low media production
                quality, bad tournament setup, limiting bracket software,
                anything really - we&apos;ll help you overcome them.<br /><br />

                We also aim to turn Auzom into a wildcard for essential
                competitive needs of both players and organizers. Customizable
                team and tournament mini-sites, advanced roster management,
                scrim booking, socializing, all sorts of fancy automation, you
                name it...<br /><br />

                But we aren&apos;t quite there yet. Right now we&apos;re just
                nearing the completion of the tournament platform, and are
                thinking through the rest of our journey. If you want to know
                more about our vision, make sure to check&nbsp;
                <a href="http://blog.auzom.gg">our blog</a>, where we&apos;ve
                outlined our long- and short-term plans.

              </div>
            </div>
          </div>
        </div>
        {!me && <div className={this.cni({ u: 'sectionSingle', d: 'cfa' })}>
          <div className={this.cni({ d: 'cfaInner' })}>
            <Link className={this.cni({ d: 'cfaSignup' })} to="/signup">
              SIGN UP
            </Link>
            <div className={this.cni({ d: 'cfaEvenif' })}>
              even if this isn&apos;t your cup of tea quite yet
            </div>
          </div>
        </div>}
      </div>
    );
  }
}
