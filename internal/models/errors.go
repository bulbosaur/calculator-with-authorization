package models

import "errors"

var (
	// ErrorCreatingDatabaseRecord - ошибка записи в БД
	ErrorCreatingDatabaseRecord = errors.New("an error occurred while writing to the database")

	// ErrorDivisionByZero - ошибка деления на ноль
	ErrorDivisionByZero = errors.New("division by zero")

	// ErrorEmptyBrackets - пустые скобочки
	ErrorEmptyBrackets = errors.New("the brackets are empty")

	// ErrorEmptyExpression - пустое выражение
	ErrorEmptyExpression = errors.New("expression is empty")

	//ErrorInvalidCharacter - запрещенные символы в выражении
	ErrorInvalidCharacter = errors.New("invalid characters in expression")

	// ErrorInvalidInput - невалидное выражение
	ErrorInvalidInput = errors.New("expression is not valid")

	// ErrorInvalidOperand - ошибка при введении операнда
	ErrorInvalidOperand = errors.New("an invalid operand")

	// ErrorInvalidRequestBody - ошибка тела запроса
	ErrorInvalidRequestBody = errors.New("invalid request body")

	// ErrorUserNotFound - пользователь не найдет
	ErrorUserNotFound = errors.New("user not found")

	// ErrorReceivingID - ошибка, которая возникает, не удается получить айди последней записи в БД
	ErrorReceivingID = errors.New("failed to get ID records in the database")

	// ErrorMissingOperand - пропущенный операнд
	ErrorMissingOperand = errors.New("missing operand")

	// ErrorUnclosedBracket - скобочки не согласованы
	ErrorUnclosedBracket = errors.New("the brackets in the expression are not consistent")
)
