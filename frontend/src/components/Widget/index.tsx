import {ReactElement, ReactNode} from 'react';
import cx from 'classnames';
import s from './widget.module.scss';

type WidgetProps = {
  title?: string | ReactNode
  className?: string
  children: ReactNode
}

export default function Widget({title, className, children}: WidgetProps): ReactElement {
  return (
    <section className={cx(s.widget, className)}>
      {title &&
      (typeof title === 'string' ? (
        <h5 className={s.title}>{title}</h5>
      ) : (
        <header className={s.title}>{title}</header>
      ))}
      <div>{children}</div>
    </section>
  );
}
