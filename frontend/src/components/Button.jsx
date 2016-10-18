import { PropTypes } from 'react';
import { Link } from 'react-router';

import Component from '../utils/Component';

export default class Button extends Component {
  static propTypes = {
    type: PropTypes.oneOf(['normal', 'important']),
    size: PropTypes.oneOf(['small', 'large']),
    disabled: PropTypes.bool,
    onClick: PropTypes.func,
    text: PropTypes.string,
    meta: PropTypes.any, // eslint-disable-line
    href: PropTypes.string,
    block: PropTypes.bool,
  };

  static defaultProps = {
    type: 'normal',
    size: 'large',
    disabled: false,
  };

  handleClick = (e) => {
    if (!this.props.disabled) {
      e.target.blur();
      this.props.onClick && this.props.onClick(e, this.props.meta);
    }
  };

  handleKeyUp = (e) => {
    if (
      !this.props.disabled && (
        e.which === 13 || // carriage return
        e.which === 32 // space
      )
    ) {
      this.props.onClick && this.props.onClick(e, this.props.meta);
    }
  };

  render() {
    const { href } = this.props;
    let Tag = 'div';
    if (href) Tag = Link;
    return (
      <Tag
        className={this.cni({
          m: [this.props.type, this.props.size, this.props.block && 'block'],
          s: this.props.disabled && 'disabled',
        })}
        tabIndex="0"
        onClick={this.handleClick}
        onKeyUp={this.handleKeyUp}
        to={href}
      >
        {this.props.text}
      </Tag>
    );
  }
}
