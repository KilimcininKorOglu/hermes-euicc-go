// Copyright (c) 2025 Kilimcinin Kör Oğlu <k@keremgok.tr>
// SPDX-License-Identifier: MIT

// Local reimplementation of the fork's notification-processing convenience API
// on the upstream damonto/euicc-go lpa.Client.

package main

import (
	"fmt"

	"github.com/damonto/euicc-go/lpa"
	sgp22 "github.com/damonto/euicc-go/v2"
)

// notificationProcessResult is the outcome of processing a single notification.
type notificationProcessResult struct {
	SequenceNumber sgp22.SequenceNumber
	Success        bool
	Error          error
	Removed        bool
}

// processNotifications sends the given notifications to their SM-DP+ servers,
// optionally removing them from the eUICC afterwards.
func processNotifications(client *lpa.Client, autoRemove, continueOnError bool, sequenceNumbers ...sgp22.SequenceNumber) ([]*notificationProcessResult, error) {
	if len(sequenceNumbers) == 0 {
		return nil, nil
	}

	results := make([]*notificationProcessResult, 0, len(sequenceNumbers))
	for _, seqNum := range sequenceNumbers {
		result := &notificationProcessResult{SequenceNumber: seqNum}
		removed, err := processSingleNotification(client, seqNum, autoRemove)
		result.Removed = removed
		if err != nil {
			result.Error = err
			results = append(results, result)
			if !continueOnError {
				return results, fmt.Errorf("failed to process notification %d: %w", seqNum, err)
			}
			continue
		}
		result.Success = true
		results = append(results, result)
	}
	return results, nil
}

// processAllNotifications processes every pending notification on the eUICC.
func processAllNotifications(client *lpa.Client, autoRemove, continueOnError bool) ([]*notificationProcessResult, error) {
	notifications, err := client.ListNotification()
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	if len(notifications) == 0 {
		return nil, nil
	}
	sequenceNumbers := make([]sgp22.SequenceNumber, len(notifications))
	for i, notif := range notifications {
		sequenceNumbers[i] = notif.SequenceNumber
	}
	return processNotifications(client, autoRemove, continueOnError, sequenceNumbers...)
}

// processSingleNotification retrieves, sends, and optionally removes one
// notification. It reports whether the notification was removed.
func processSingleNotification(client *lpa.Client, seqNum sgp22.SequenceNumber, autoRemove bool) (bool, error) {
	notifications, err := client.RetrieveNotificationList(seqNum)
	if err != nil {
		return false, fmt.Errorf("retrieve notification: %w", err)
	}
	if len(notifications) == 0 {
		return false, fmt.Errorf("notification with sequence number %d not found", seqNum)
	}

	if err := client.HandleNotification(notifications[0]); err != nil {
		return false, fmt.Errorf("handle notification: %w", err)
	}

	if autoRemove {
		if err := client.RemoveNotificationFromList(seqNum); err != nil {
			return false, fmt.Errorf("remove notification: %w", err)
		}
		return true, nil
	}
	return false, nil
}
