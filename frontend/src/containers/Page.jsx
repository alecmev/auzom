import { PropTypes } from 'react';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { Link, IndexLink } from 'react-router';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import Component from '../utils/Component';
import Footer from './Footer';

import Img from '../components/Img';

const messageSelector = createSelector(
  selectors.messages,
  messages => messages.first(),
);

@connect(createStructuredSelector({
  myId: selectors.myId,
  amAdmin: selectors.amAdmin,
  me: selectors.me,
  isLoading: selectors.isLoading,
  canGoBack: selectors.canGoBack,
  message: messageSelector,
}), actions)
export default class Page extends Component {
  static propTypes = {
    myId: PropTypes.string,
    amAdmin: PropTypes.bool,
    me: ImmutablePropTypes.map,
    isLoading: PropTypes.bool.isRequired,
    canGoBack: PropTypes.bool.isRequired,
    message: ImmutablePropTypes.map,
    children: PropTypes.element,
    location: PropTypes.object.isRequired, // eslint-disable-line
    loadUser: PropTypes.func.isRequired,
    locationChanged: PropTypes.func.isRequired,
    logOut: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  state = {
    isScrolled: !!window.scrollY,
    isExpanded: false,
    scrollX: window.scrollX,
  };

  componentWillMount() {
    if (this.props.myId) {
      this.props.loadUser(this.props.myId);
    }

    // There's no need in resetting state, this component instance is gone once
    // unmounted anyway. https://github.com/facebook/react/issues/4345
  }

  componentDidMount() {
    window.addEventListener('scroll', this.handleScroll);
  }

  componentWillReceiveProps(nextProps) {
    if (!nextProps.myId) {
      this.setState({ isExpanded: false });
    } else if (nextProps.myId !== this.props.myId) {
      this.props.loadUser(nextProps.myId);
    }

    if (nextProps.location !== this.props.location) {
      if (nextProps.myId) {
        this.setState({ isExpanded: false });
      }

      if (!this.props.canGoBack) {
        this.props.locationChanged();
      }
    }
  }

  componentWillUnmount() {
    window.removeEventListener('scroll', this.handleScroll);
  }

  handleScroll = () => {
    if (this.state.isScrolled !== !!window.scrollY) {
      this.setState({ isScrolled: !!window.scrollY });
    }
    if (this.state.scrollX !== window.scrollX) {
      this.setState({ scrollX: window.scrollX });
    }
  };

  handleGravatar = () =>
    this.setState({ isExpanded: !this.state.isExpanded });

  handleLogout = () =>
    this.props.logOut(() => this.props.messagePush('logged out'));

  render() {
    const { isExpanded, isScrolled } = this.state;
    const { message: _message, me: _me } = this.props;

    let Message = null;
    if (_message) {
      const message = _message.toJS();
      const level = message.isError ? 'error' : 'success';
      Message = (
        <div className={this.cn({
          u: 'sectionSingle',
          d: 'message',
          m: level,
        })}>
          <div className={this.cn({ u: 'sectionMargined' })}>
            <span className={this.cn({ d: 'messageLevel', m: level })}>
              {level.toUpperCase()}
              {message.code && ` ${message.code}`}
            </span>
            {message.url && <span className={this.cn({ d: 'messageURL' })}>
              {message.url}
            </span>}
            {message.text[0].toUpperCase() + message.text.slice(1)}
          </div>
        </div>
      );
    }

    let me = null;
    if (_me) me = _me.toJS();

    return (
      <div className={this.cni()}>
        <div className={this.cn({
          u: 'sectionSingle',
          d: 'header',
          s: [isScrolled && 'scrolled', isExpanded && 'expanded'],
        })} style={{ left: -this.state.scrollX }}>
          <div className={this.cn({ u: 'sectionMargined' })}>
            <div className={this.cn({ d: 'menu' })}>
              <div className={this.cn({ d: 'menuLogo' })}>
                <Link to="/">
                  <img
                    className={this.cn({
                      d: 'menuLogoImage',
                      s: this.props.isLoading && 'loading',
                    })}
                    src={require('../assets/logo.svg')}
                  />
                  <span className={this.cn({ d: 'menuLogoTitle' })}>auzom</span>
                </Link>
                <Link
                  to="/about"
                  className={this.cn({ d: 'menuLogoSubtitle' })}
                >
                  α
                  <div className={this.cn({ d: 'menuLogoSubtitleTooltip' })}>
                    alpha
                  </div>
                </Link>
              </div>
              <Link
                className={this.cn({ d: 'menuEntry' })}
                activeClassName="is-active"
                to="/battlefield-4/conquest-league"
              >
                bcl
              </Link>
              <Link
                className={this.cn({ d: 'menuEntry' })}
                activeClassName="is-active"
                to="/battlefield-4/aces-high-ah"
              >
                ace
              </Link>
              <IndexLink
                className={this.cn({ d: 'menuEntry' })}
                activeClassName="is-active"
                to="/teams"
              >
                teams
              </IndexLink>
              {me &&
                <Link
                  className={this.cn({ d: 'menuEntry' })}
                  activeClassName="is-active"
                  to="/teams/new"
                >
                  new team
                </Link>
              }
              <div className={this.cn({ d: 'menuFiller' })} />
              {/* <div className={this.cn({ d: 'menuNotifications' })}>
                <span className="octicon octicon-bell"></span>
              </div> */}
              {me ? [
                <div
                  className={this.cn({
                    d: 'menuProfile',
                    s: isExpanded && 'expanded',
                  })}
                  key="profile"
                >
                  <Link
                    className={this.cn({ d: 'menuProfileName' })}
                    to={`/users/${me.id}`}
                  >
                    {me.nickname}
                  </Link>
                  <div className={this.cn({ d: 'menuProfileEmail' })}>
                    {me.email}
                  </div>
                </div>,
                <Img
                  className={this.cn({
                    d: 'menuGravatar',
                    s: isExpanded && 'expanded',
                  })}
                  src={`//gravatar.com/avatar/${me.gravatar}?s=96&d=identicon`}
                  onClick={this.handleGravatar}
                  key="gravatar"
                />,
              ] : [
                <Link
                  className={this.cn({ d: 'menuEntry' })}
                  activeClassName="is-active"
                  to="/login"
                  key="login"
                >
                  log in
                </Link>,
                <Link
                  className={this.cn({ d: 'menuEntry', m: 'signup' })}
                  activeClassName="is-active"
                  to="/signup"
                  key="signup"
                >
                  sign up
                </Link>,
              ]}
            </div>
            <div className={this.cn({
              d: 'fancystuff',
              s: isExpanded && 'expanded',
            })}>
              there will be some fancy stuff here later on...
            </div>
            <div className={this.cn({
              d: 'profile',
              s: isExpanded && 'expanded',
            })}>
              <Link
                className={this.cn({ d: 'profileButton' })}
                to="/settings"
              >
                settings
              </Link><br />
              <button
                className={this.cn({ d: 'profileButton' })}
                onClick={this.handleLogout}
              >
                log out
              </button>
            </div>
            <div className={this.cn({ u: 'clear' })} />
          </div>
          <ReactCSSTransitionGroup
            transitionName={this.cn({ d: 'messageTransition' })}
            transitionEnterTimeout={250}
            transitionLeaveTimeout={250}
          >
            {Message}
          </ReactCSSTransitionGroup>
        </div>
        <div className={this.cn({
          d: 'theRest',
          s: isExpanded && 'muted',
        })}>
          <div className={this.cn({ d: 'theRestInner' })}>
            {this.props.children}
            <Footer />
          </div>
        </div>
      </div>
    );
  }
}
