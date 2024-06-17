package transactionmanager

/** Примеры вложенных транзакций
 * псевдо-вложенная транзакция
 * вложенная субтранзакция
 * вложенная независимая транзакция

// Примитивный пример
func foo() (err error) {
	defer func() {
		// вложенная независимая транзакция
		// do smth
		if err != nil {
			logger.Error(err)
		}
	}()

	err = func() error {
		// псевдо-вложенные транзакции
		// do smth
		return nil
	}()
	if err != nil {
		return err
	}

	go func() {
		// вложенная субтранзакция
		// do smth
		if err != nil {
			logger.Error(err)
		}
	}()

	return nil
}

type service struct {
	trm trm.TransactionManager
}

// Пример с менеджером транзакций
func (s service) foo(ctx context.Context) error {
	/// Вложенная независимая транзакция
	// Не зависит от исхода выполнения основной транзакции и
	// сама не влияет на ее исход
	// Выполняется в любом случае. Фиксируется, если сама завершается успешно
	defer func() {
		err := s.trm.DoWithConfig(ctx,
			config.MustTrConfig(trm.NestedIndependentTr),
			func(ctx context.Context) error {
				// Независимая транзакция
				// ...
				return err
			},
		)
		if err != nil {
			logger.Error("message", err)
		}
	}()

	/// Псевдо-вложенная транзакция
	// Зависит от исхода внешних и внутренних псевдо-вложенных транзакций.
	// Ошибка приводит к откату всех внешних псевдо-вложенных транзакций и
	// субтранзакций
	// Фиксируется, если все внешние и все внутренние псевдо-вложенные транзакции
	// завершаются успехом и сама завершается успехом
	err := s.trm.DoWithConfig(ctx,
		config.MustTrConfig(trm.PseudoNestedTr),
		func(ctx context.Context) error {
			// Псевдо-вложенная транзакция
			// ...
			return err
		},
	)
	if err != nil {
		logger.Error("message", err)
		return err
	}

	/// Вложенная субтранзакция (дополнительная)
	// Зависит от исхода выполнения основной транзакции, но сама
	// не влияет на ее исход.
	// Выполняется, только если основная транзакция завершилась успешно.
	// Фиксируется, если фиксируется основная транзакция и сама завершается успешно
	go func() {
		s.trm.DoWithConfig(ctx,
			config.MustTrConfig(trm.NestedSubTr),
			func(ctx context.Context) error {
				// Субранзакция
				// ...
				if err != nil {
					logger.Error("message", err)
				}

				return nil
			},
		)

	}()

	return nil
}

**/
